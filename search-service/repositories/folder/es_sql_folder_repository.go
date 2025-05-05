package folder

import (
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

func NewEsSqlFolderRepository(sqlDb *infrastructure.GormAdapter, esDb *infrastructure.ElasticAdapter) *EsSqlFolderRepository {
	return &EsSqlFolderRepository{sql: sqlDb, es: esDb}
}

// CreateFolder implements FolderRepository.
func (e *EsSqlFolderRepository) CreateFolder(query CreateFolderQuery) (*CreateFolderResult, error) {
	createdFolder := types.Folder{
		ID:       query.ID,
		Name:     query.Name,
		ParentID: &query.ParentID,
	}
	if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(createdFolder).Error; err != nil {
			return err
		}
		// TO DO: Create the folder in Elasticsearch
		return nil
	}); err != nil {
		return nil, err
	}
	return &CreateFolderResult{Folder: createdFolder}, nil
}

// DeleteFolder implements FolderRepository.
func (e *EsSqlFolderRepository) DeleteFolder(query DeleteFolderQuery) error {
	if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&types.Folder{}, query.ID).Error; err != nil {
			return err
		}
		// TO DO: Delete the folder in Elasticsearch
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetFolder implements FolderRepository.
func (e *EsSqlFolderRepository) GetFolder(query GetFolderQuery) (*GetFolderResult, error) {
	// TODO: Implement the logic to retrieve a folder from Elasticsearch
	panic("unimplemented")
}

// GetFolderHierarchy implements FolderRepository.
func (e *EsSqlFolderRepository) GetFolderHierarchy(query GetFolderHierarchyQuery) (*GetFolderHierarchyResult, error) {
	var folders []types.Folder
	request := `
        WITH RECURSIVE folder_hierarchy AS (
            SELECT id, name, parent_id
            FROM folders
            WHERE id = ?
            UNION ALL
            SELECT f.id, f.name, f.parent_id
            FROM folders f
            INNER JOIN folder_hierarchy fh ON fh.id = f.parent_id
        )
        SELECT * FROM folder_hierarchy;
    `

	if err := e.sql.Db.Raw(request, query.ID).Scan(&folders).Error; err != nil {
		return nil, err
	}

	return &GetFolderHierarchyResult{Folders: folders}, nil
}

// SearchFolder implements FolderRepository.
func (e *EsSqlFolderRepository) SearchFolder(query SearchFolderQuery) (*SearchFolderResult, error) {
	// TO DO: Search the folder in Elasticsearch
	panic("unimplemented")
}

// UpdateFolder implements FolderRepository.
func (e *EsSqlFolderRepository) UpdateFolder(query UpdateFolderQuery) (*UpdateFolderResult, error) {
	res, err := e.GetFolder(GetFolderQuery{ID: query.ID})
	if err != nil {
		return nil, err
	}
	if err := e.sql.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&types.Folder{}).Where("id = ?", query.ID).Updates(types.Folder{
			Name: query.Name,
		}).Error; err != nil {
			return err
		}
		// TO DO: Update the folder in Elasticsearch
		return nil
	}); err != nil {
		return nil, err
	}
	return &UpdateFolderResult{Folder: types.Folder{
		ID:       query.ID,
		Name:     query.Name,
		ParentID: res.Folder.ParentID,
	}}, nil
}
