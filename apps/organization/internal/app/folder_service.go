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

func (s *FolderService) UpdateFolder(folder organization.Folder) error {

	data := avroGeneratedSchema.FolderEvent{
		Id:          folder.Id,
		Name:        folder.Name,
		Description: StringPtrToValue(folder.Description),
		Icon:        StringPtrToValue(folder.Icon),
		Created_at:  folder.CreatedAt.String(),
		Updated_at:  folder.UpdatedAt.String(),
		Parent_id:   StringPtrToValue(folder.ParentID),
		Members:     folder.Members,
		Created_by:  folder.CreatedBy,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}
	return s.publisher.Publish("Folder-Update", buf.Bytes())
}

func (s *FolderService) DeleteFolder(id string) error {
	data := map[string]interface{}{
		"id": id,
	}
	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}
	return s.publisher.Publish("Folder-Delete", encoded)
}

func (s *FolderService) ListFolders(req organization.GetFolderRequest) ([]organization.Folder, error) {
	var folders []organization.Folder
	if err := s.database.Model(&organization.Folder{}).Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("created_at asc").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (s *FolderService) GetFolder(id string) (organization.Folder, error) {
	// TODO: Replace with real implementation
	return organization.Folder{Id: id}, nil
}
