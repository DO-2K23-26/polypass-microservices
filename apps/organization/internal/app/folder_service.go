package app

import (
	"bytes"
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
	newFolder := organization.Folder{
		Id:          uuid.New().String(),
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ParentID:    folder.ParentID,
		Members:     []string{folder.CreatedBy},
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
		Members:     newFolder.Members,
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

	err := s.publisher.Publish("Folder-Create", buf.Bytes())
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

	updatedFolder := organization.Folder{
		Id:          previousFolder.Id,
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   previousFolder.CreatedAt,
		UpdatedAt:   time.Now(),
		ParentID:    folder.ParentID,
		Members:     []string{previousFolder.CreatedBy},
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
		Members:     updatedFolder.Members,
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

	kafkaErr := s.publisher.Publish("Folder-Update", buf.Bytes())
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

	s.database.Delete(&organization.Folder{}, "id = ?", folderId)
	data := avroGeneratedSchema.FolderEvent{
		Id:          deletedFolder.Id,
		Name:        deletedFolder.Name,
		Description: StringPtrToValue(deletedFolder.Description),
		Icon:        StringPtrToValue(deletedFolder.Icon),
		Created_at:  deletedFolder.CreatedAt.String(),
		Updated_at:  deletedFolder.UpdatedAt.String(),
		Parent_id:   StringPtrToValue(deletedFolder.ParentID),
		Members:     deletedFolder.Members,
		Created_by:  deletedFolder.CreatedBy,
	}
	var buf bytes.Buffer
	encodeErr := data.Serialize(&buf)
	if encodeErr != nil {
		return encodeErr
	}

	return s.publisher.Publish("Folder-Delete", buf.Bytes())
}

func (s *FolderService) ListFolders(req organization.GetFolderRequest) ([]organization.Folder, error) {
	var folders []organization.Folder
	if err := s.database.Model(&organization.Folder{}).Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("created_at asc").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (s *FolderService) GetFolder(id string) (organization.Folder, error) {
	var folder organization.Folder
	if err := s.database.Find(&folder, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return organization.Folder{}, gorm.ErrRecordNotFound
		}
		return organization.Folder{}, err
	}
	return folder, nil
}
