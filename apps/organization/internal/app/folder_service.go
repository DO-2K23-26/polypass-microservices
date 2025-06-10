package app

import (
	"bytes"
	"time"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"github.com/google/uuid"
)

type EventPublisher interface {
	Publish(topic string, data []byte) error
}

type FolderService struct {
	publisher EventPublisher
	encoder   *schemautils.AvroEncoder
}

func NewFolderService(publisher EventPublisher, encoder *schemautils.AvroEncoder) *FolderService {
	return &FolderService{publisher: publisher, encoder: encoder}
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

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
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
