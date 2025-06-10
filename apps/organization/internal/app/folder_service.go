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

func (s *FolderService) CreateFolder(folder organization.CreateFolderRequest) error {
	data := avroGeneratedSchema.FolderEvent{
		Id:          uuid.New().String(),
		Name:        folder.Name,
		Description: StringPtrToValue(folder.Description),
		Icon:        StringPtrToValue(folder.Icon),
		Created_at:  time.Now().String(),
		Updated_at:  time.Now().String(),
		Parent_id:   StringPtrToValue(folder.ParentID),
		Members:     []string{},
		Created_by:  folder.CreatedBy,
	}

	res := s.database.Create(&organization.Folder{
		Id:          data.Id,
		Name:        folder.Name,
		Description: folder.Description,
		Icon:        folder.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ParentID:    folder.ParentID,
		Members:     data.Members,
		CreatedBy:   folder.CreatedBy,
	})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	fmt.Println("Folder created with ID:", data.Id)

	var buf bytes.Buffer
	kafkaErr := data.Serialize(&buf)
	if kafkaErr != nil {
		return kafkaErr
	}

	return s.publisher.Publish("Folder-Create", buf.Bytes())
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

func (s *FolderService) ListFolders() ([]organization.Folder, error) {
	// TODO: Replace with real implementation
	return []organization.Folder{}, nil
}

func (s *FolderService) GetFolder(id string) (organization.Folder, error) {
	// TODO: Replace with real implementation
	return organization.Folder{Id: id}, nil
}
