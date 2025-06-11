package folder

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	esTypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"

	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"gorm.io/gorm"
)

// This struct represents a repository for managing folders using SQL and Elasticsearch.
// We will duplicate data in Elasticsearch and PostgreSQL.
type EsSqlFolderRepository struct {
	sql *infrastructure.GormAdapter
	es  *infrastructure.ElasticAdapter
}

func NewEsSqlFolderRepository(sqlDb *infrastructure.GormAdapter, esDb *infrastructure.ElasticAdapter) IFolderRepository {
	return &EsSqlFolderRepository{sql: sqlDb, es: esDb}
}

// GetFromUser implements IFolderRepository.
func (e *EsSqlFolderRepository) GetFromUser(query GetFromUserQuery) (*GetFromUserResult, error) {
	var user types.User
	if err := e.sql.Db.Preload("Folders").Where("id = ?", query.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	folderIDs := make([]string, len(user.Folders))
	for i, folder := range user.Folders {
		folderIDs[i] = folder.ID
	}

	hierarchyQuery := MGetFolderHierarchyQuery{
		IDs: folderIDs,
	}

	hierarchyResult, err := e.MGetHierarchy(hierarchyQuery)
	if err != nil {
		return nil, err
	}

	return &GetFromUserResult{
		Folders: hierarchyResult.Folders,
	}, nil
}

// CreateFolder implements FolderRepository.
func (e *EsSqlFolderRepository) Create(query CreateFolderQuery) (*CreateFolderResult, error) {
	// Vérifier si le dossier parent existe si un parent_id est spécifié
	if query.ParentID != "" {
		var parentFolder types.Folder
		if err := e.sql.Db.Where("id = ?", query.ParentID).First(&parentFolder).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent folder with ID %s does not exist", query.ParentID)
			}
			return nil, err
		}
	}

	createdFolder := types.Folder{
		ID:          query.ID,
		Name:        query.Name,
		Description: query.Description,
		Icon:        query.Icon,
		ParentID:    query.ParentID,
		CreatedBy:   query.CreatedBy,
		CreatedAt:   "", // Ces champs seront remplis par la base de données
		UpdatedAt:   "", // Ces champs seront remplis par la base de données
	}
	createdFolder.SetMembers(query.Members)

	if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&createdFolder).Error; err != nil {
			return err
		}
		log.Println("create folder", createdFolder)
		if err := e.es.CreateDocument(types.FolderIndex, createdFolder.ID, &createdFolder); err != nil {
			log.Println("create folder err: ", err)
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &CreateFolderResult{Folder: createdFolder}, nil
}

// DeleteFolder implements FolderRepository.
func (e *EsSqlFolderRepository) Delete(query DeleteFolderQuery) error {
	res, err := e.GetHierarchy(GetFolderHierarchyQuery{ID: query.ID})
	if err != nil {
		return err
	}

	for _, folder := range res.Folders {
		if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id", folder.ID).Delete(&types.Folder{}).Error; err != nil {
				return err
			}
			e.es.DeleteByQuery(types.CredentialIndex, esTypes.Query{
				Term: map[string]esTypes.TermQuery{
					"folder.id": {
						Value: folder.ID,
					},
				},
			})
			e.es.DeleteByQuery(types.TagIndex, esTypes.Query{
				Term: map[string]esTypes.TermQuery{
					"folder_id": {
						Value: folder.ID,
					},
				},
			})
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// GetFolder implements FolderRepository.
func (e *EsSqlFolderRepository) Get(query GetFolderQuery) (*GetFolderResult, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}
	var folder types.Folder

	err := e.es.GetDocument(types.TagIndex, query.ID, &folder)
	if err != nil {
		return nil, fmt.Errorf("Error getting tag: %w", err)
	}

	return &GetFolderResult{
		Folder: folder,
	}, nil
}

// GetFolderHierarchy implements FolderRepository.
func (e *EsSqlFolderRepository) GetHierarchy(query GetFolderHierarchyQuery) (*GetFolderHierarchyResult, error) {
	var folders []types.Folder
	request := `
				WITH RECURSIVE folder_hierarchy AS (
								SELECT id, name, parent_id, 0 AS depth
								FROM folders
								WHERE id = ?
								UNION ALL
								SELECT f.id, f.name, f.parent_id, fh.depth + 1 AS depth
								FROM folders f
								INNER JOIN folder_hierarchy fh ON fh.id = f.parent_id
				)
				SELECT id, name, parent_id
				FROM folder_hierarchy
				ORDER BY depth DESC;
				`

	if err := e.sql.Db.Raw(request, query.ID).Scan(&folders).Error; err != nil {
		return nil, err
	}

	return &GetFolderHierarchyResult{Folders: folders}, nil
}

// GetFolderHierarchy implements FolderRepository.
func (e *EsSqlFolderRepository) MGetHierarchy(query MGetFolderHierarchyQuery) (*MGetFolderHierarchyResult, error) {
	var folders []types.Folder
	request := `
				WITH RECURSIVE folder_hierarchy AS (
								SELECT id, name, parent_id, 0 AS depth
								FROM folders
								WHERE id IN (?)
								UNION ALL
								SELECT f.id, f.name, f.parent_id, fh.depth + 1 AS depth
								FROM folders f
								INNER JOIN folder_hierarchy fh ON fh.id = f.parent_id
				)
				SELECT id, name, parent_id, depth
				FROM (
					SELECT id, name, parent_id, depth, ROW_NUMBER() OVER (PARTITION BY id ORDER BY depth DESC) AS row_num
					FROM folder_hierarchy
				) AS ranked_folders
				WHERE row_num = 1
				ORDER BY depth DESC;
				`

	if err := e.sql.Db.Raw(request, query.IDs).Scan(&folders).Error; err != nil {
		return nil, err
	}

	return &MGetFolderHierarchyResult{Folders: folders}, nil
}

// SearchFolder implements FolderRepository.
func (e *EsSqlFolderRepository) Search(query SearchFolderQuery) (*SearchFolderResult, error) {
	// Default limit and offset if not provided
	limit := 10
	if query.Limit != nil {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}

	// Construct the search query
	res, total, err := e.es.Search(
		types.FolderIndex,
		query.Name,
		[]string{"name", "description"},
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Parse the search results
	folders := make([]types.Folder, *total)
	for i, hit := range res {
		if err := json.Unmarshal(hit, &folders[i]); err != nil {
			return nil, fmt.Errorf("error unmarshalling hit source: %w", err)
		}
	}

	return &SearchFolderResult{
		Folders: folders,
		Total:   *total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

// UpdateFolder implements FolderRepository.
func (e *EsSqlFolderRepository) Update(query UpdateFolderQuery) (*UpdateFolderResult, error) {
	params := map[string]json.RawMessage{}
	val, _ := json.Marshal(query.ID)
	params["folder_id"] = val

	// Marshal values
	val, _ = json.Marshal(query.Name)
	params["folder_name"] = val

	val, _ = json.Marshal(query.Description)
	params["description"] = val

	val, _ = json.Marshal(query.Icon)
	params["icon"] = val

	val, _ = json.Marshal(query.ParentID)
	params["parent_id"] = val

	val, _ = json.Marshal(query.Members)
	params["members"] = val

	script := `
if (ctx._source.folder != null) {
  if (params.containsKey('folder_name')) {
    ctx._source.folder.name = params.folder_name;
  }
  if (params.containsKey('description')) {
    ctx._source.folder.description = params.description;
  }
  if (params.containsKey('icon')) {
    ctx._source.folder.icon = params.icon;
  }
  if (params.containsKey('parent_id')) {
    ctx._source.folder.parent_id = params.parent_id;
  }
  if (params.containsKey('members')) {
    ctx._source.folder.members = params.members;
  }
}`

	esQuery := &esTypes.Query{
		Term: map[string]esTypes.TermQuery{
			"folder.id": {
				Value: query.ID,
			},
		},
	}

	// Créer une structure de mise à jour avec le champ Members converti en string
	updateData := types.Folder{
		ID:          query.ID,
		Name:        query.Name,
		Description: query.Description,
		Icon:        query.Icon,
		ParentID:    query.ParentID,
	}
	updateData.SetMembers(query.Members)

	if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&types.Folder{}).Where("id = ?", query.ID).Updates(updateData).Error; err != nil {
			return err
		}

		err := e.es.UpdateDocument(types.FolderIndex, query.ID, updateData)
		if err != nil {
			return err
		}

		err = e.es.UpdateByQuery(types.CredentialIndex, *esQuery, script, &params)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Récupérer le dossier mis à jour pour le retourner
	var updatedFolder types.Folder
	if err := e.sql.Db.Where("id = ?", query.ID).First(&updatedFolder).Error; err != nil {
		return nil, err
	}

	return &UpdateFolderResult{Folder: updatedFolder}, nil
}
