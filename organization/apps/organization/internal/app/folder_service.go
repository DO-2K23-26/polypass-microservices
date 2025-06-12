package app

import (
	"bytes"
	"fmt"
	"time"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventPublisher interface {
	Publish(topic string, data []byte) error
}

type FolderService struct {
	publisher EventPublisher
	encoder   *schemautils.AvroEncoder
	database  *gorm.DB
}

func NewFolderService(publisher EventPublisher, encoder *schemautils.AvroEncoder, database *gorm.DB) *FolderService {
	return &FolderService{publisher: publisher, encoder: encoder, database: database}
}

func StringPtrToValue(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func (s *FolderService) CreateFolder(folder organization.CreateFolderRequest) (*organization.Folder, error) {
	if folder.ParentID != nil {
		_, err := s.GetFolder(*folder.ParentID)
		if err != nil {
			return nil, err
		}
	}

	newFolder := organization.Folder{
		Id:          uuid.New().String(),
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ParentID:    folder.ParentID,
		User:        &[]organization.User{{ID: folder.CreatedBy}},
		CreatedBy:   folder.CreatedBy,
	}

	data := avroGeneratedSchema.FolderEvent{
		Id:          newFolder.Id,
		Name:        newFolder.Name,
		Description: StringPtrToValue(newFolder.Description),
		Icon:        StringPtrToValue(newFolder.Icon),
		Created_at:  newFolder.CreatedAt.String(),
		Updated_at:  newFolder.UpdatedAt.String(),
		Parent_id:   StringPtrToValue(newFolder.ParentID),
		Members:     []string{folder.CreatedBy},
		Created_by:  newFolder.CreatedBy,
	}

	res := s.database.Create(&newFolder)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var buf bytes.Buffer
	kafkaErr := data.Serialize(&buf)
	if kafkaErr != nil {
		return nil, kafkaErr
	}

	err := s.publisher.Publish("folder-creation", buf.Bytes())
	if err != nil {
		return nil, err
	}

	return &newFolder, nil
}

func (s *FolderService) UpdateFolder(folderId string, folder organization.UpdateFolderRequest) (*organization.Folder, error) {
	previousFolder, getDatabaseErr := s.GetFolder(folderId)
	if getDatabaseErr != nil {
		if getDatabaseErr == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, getDatabaseErr
	}

	if folder.ParentID != nil {
		_, err := s.GetFolder(*folder.ParentID)
		if err != nil {
			return nil, err
		}
	}

	updatedFolder := organization.Folder{
		Id:          previousFolder.Id,
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   previousFolder.CreatedAt,
		UpdatedAt:   time.Now(),
		ParentID:    folder.ParentID,
		User:        &[]organization.User{{ID: previousFolder.CreatedBy}},
		CreatedBy:   previousFolder.CreatedBy,
	}

	data := avroGeneratedSchema.FolderEvent{
		Id:          updatedFolder.Id,
		Name:        updatedFolder.Name,
		Description: StringPtrToValue(updatedFolder.Description),
		Icon:        StringPtrToValue(updatedFolder.Icon),
		Created_at:  updatedFolder.CreatedAt.String(),
		Updated_at:  updatedFolder.UpdatedAt.String(),
		Parent_id:   StringPtrToValue(updatedFolder.ParentID),
		Members:     []string{previousFolder.CreatedBy},
		Created_by:  updatedFolder.CreatedBy,
	}

	updateDatabase := s.database.Model(&organization.Folder{}).Where("id = ?", folderId).Updates(updatedFolder)
	if updateDatabase.Error != nil {
		return nil, updateDatabase.Error
	}

	var buf bytes.Buffer
	serializeErr := data.Serialize(&buf)
	if serializeErr != nil {
		return nil, serializeErr
	}

	kafkaErr := s.publisher.Publish("folder-update", buf.Bytes())
	if kafkaErr != nil {
		return nil, kafkaErr
	}

	return &updatedFolder, nil
}

func (s *FolderService) DeleteFolder(folderId string) error {
	deletedFolder, getDatabaseErr := s.GetFolder(folderId)
	if getDatabaseErr != nil {
		if getDatabaseErr == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return getDatabaseErr
	}

	// Delete all users associated with the folder
	if err := s.database.Model(deletedFolder).
		Association("User").Clear(); err != nil {
		return fmt.Errorf("failed to clear users from folder: %w", err)
	}

	err := s.database.Delete(&organization.Folder{}, "id = ?", folderId).Error
	if err != nil {
		return err
	}

	data := avroGeneratedSchema.FolderEvent{
		Id:          deletedFolder.Id,
		Name:        deletedFolder.Name,
		Description: StringPtrToValue(deletedFolder.Description),
		Icon:        StringPtrToValue(deletedFolder.Icon),
		Created_at:  deletedFolder.CreatedAt.String(),
		Updated_at:  deletedFolder.UpdatedAt.String(),
		Parent_id:   StringPtrToValue(deletedFolder.ParentID),
		Members:     []string{deletedFolder.CreatedBy},
		Created_by:  deletedFolder.CreatedBy,
	}
	var buf bytes.Buffer
	encodeErr := data.Serialize(&buf)
	if encodeErr != nil {
		return encodeErr
	}

	return s.publisher.Publish("folder-delete", buf.Bytes())
}

func (s *FolderService) ListFolders(req organization.GetFolderRequest) ([]organization.Folder, error) {
	var folders []organization.Folder

	if req.UserId != nil {
		// Retrieve folders that the user is a direct member of
		if err := s.database.
			Joins("JOIN user_folders ON user_folders.folder_id = folders.id").
			Where("user_folders.user_id = ?", *req.UserId).
			Find(&folders).Error; err != nil {
			return nil, err
		}

		// Use a queue to manage folders whose children need to be retrieved
		var queue = folders
		for len(queue) > 0 {
			var nextQueue []organization.Folder
			for _, folder := range queue {
				// Retrieve children folders
				var childFolders []organization.Folder
				if err := s.database.
					Where("parent_id = ?", folder.Id).
					Find(&childFolders).Error; err != nil {
					return nil, err
				}

				// Append child folders to the result
				folders = append(folders, childFolders...)
				nextQueue = append(nextQueue, childFolders...)
			}
			queue = nextQueue
		}

		// Remove duplicates from the folders slice
		folderMap := make(map[string]organization.Folder)
		for _, folder := range folders {
			if _, exists := folderMap[folder.Id]; !exists {
				folderMap[folder.Id] = folder
			}
		}
		folders = make([]organization.Folder, 0, len(folderMap))
		for _, folder := range folderMap {
			folders = append(folders, folder)
		}

		// Limit and paginate the results
		if len(folders) > req.Limit {
			if req.Page < 1 {
				req.Page = 1
			}
			if req.Page > (len(folders)+req.Limit-1)/req.Limit {
				req.Page = (len(folders) + req.Limit - 1) / req.Limit
			}
			start := (req.Page - 1) * req.Limit
			end := start + req.Limit
			if end > len(folders) {
				end = len(folders)
			}
			folders = folders[start:end]
		}

		// Sort folders by created_at in ascending order without using gorm
		if len(folders) > 0 {
			for i := 0; i < len(folders)-1; i++ {
				for j := i + 1; j < len(folders); j++ {
					if folders[i].CreatedAt.After(folders[j].CreatedAt) {
						folders[i], folders[j] = folders[j], folders[i]
					}
				}
			}
		}
	} else {
		// Retrieve all folders without user filter
		if err := s.database.
			Limit(req.Limit).
			Offset((req.Page - 1) * req.Limit).
			Order("created_at asc").
			Find(&folders).Error; err != nil {
			return nil, err
		}
	}

	return folders, nil
}

func (s *FolderService) GetFolder(id string) (*organization.Folder, error) {
	var folder organization.Folder
	if err := s.database.Find(&folder, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if folder.Id == "" {
		return nil, fmt.Errorf("folder with id %s not found", id)
	}

	return &folder, nil
}

func (s *FolderService) ListUsersInFolder(folderId string, req organization.GetUsersInFolderRequest) ([]string, error) {
	// Retrieve all parent folders up to the root
	var folders []organization.Folder
	currentFolderId := folderId
	for {
		var f organization.Folder
		result := s.database.Preload("User").First(&f, "id = ?", currentFolderId)
		if result.Error != nil {
			return nil, result.Error
		}
		folders = append(folders, f)
		if f.ParentID == nil || *f.ParentID == "" {
			break
		}
		currentFolderId = *f.ParentID
	}

	// Collect all user IDs without duplicates
	userIDSet := make(map[string]struct{})
	for _, f := range folders {
		if f.User != nil {
			for _, user := range *f.User {
				userIDSet[user.ID] = struct{}{}
			}
		}
	}

	var userIDs []string
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
