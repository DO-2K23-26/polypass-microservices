package app

import (
	"bytes"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
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

func (s *FolderService) CreateFolder(folder organization.Folder) error {
	data := avroGeneratedSchema.FolderEvent{
		Id:          folder.Id,
		Name:        folder.Name,
		Description: &avroGeneratedSchema.UnionNullString{String: *folder.Description},
		Icon:        &avroGeneratedSchema.UnionNullString{String: *folder.Icon},
		Created_at:  folder.CreatedAt.String(),
		Updated_at:  folder.UpdatedAt.String(),
		Parent_id:   &avroGeneratedSchema.UnionNullString{String: *folder.ParentID},
		Members:     folder.Members,
		Created_by:  folder.CreatedBy,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	// encoded, err := s.encoder.Encode(data)
	// if err != nil {
	// 	return err
	// }

	return s.publisher.Publish("Folder-Create", buf.Bytes())
}

func (s *FolderService) UpdateFolder(folder organization.Folder) error {
	data := avroGeneratedSchema.FolderEvent{
		Id:          folder.Id,
		Name:        folder.Name,
		Description: &avroGeneratedSchema.UnionNullString{String: *folder.Description},
		Icon:        &avroGeneratedSchema.UnionNullString{String: *folder.Icon},
		Created_at:  folder.CreatedAt.String(),
		Updated_at:  folder.UpdatedAt.String(),
		Parent_id:   &avroGeneratedSchema.UnionNullString{String: *folder.ParentID},
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
