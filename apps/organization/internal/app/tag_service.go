package app

import (
	"bytes"

	avroGeneratedSchema "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
)

type TagService struct {
	publisher EventPublisher
	encoder   *schemautils.AvroEncoder
}

func NewTagService(publisher EventPublisher, encoder *schemautils.AvroEncoder) *TagService {
	return &TagService{publisher: publisher, encoder: encoder}
}

func (s *TagService) CreateTag(tag organization.Tag) error {
	data := avroGeneratedSchema.TagEvent{
		Id:         tag.Id,
		Name:       tag.Name,
		Color:      tag.Color,
		Created_at: tag.CreatedAt.String(),
		Updated_at: tag.UpdatedAt.String(),
		Folder_id:  tag.FolderID,
		Created_by: tag.CreatedBy,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	return s.publisher.Publish("Tag-Create", buf.Bytes())
}

func (s *TagService) UpdateTag(tag organization.Tag) error {
	data := avroGeneratedSchema.TagEvent{
		Id:         tag.Id,
		Name:       tag.Name,
		Color:      tag.Color,
		Created_at: tag.CreatedAt.String(),
		Updated_at: tag.UpdatedAt.String(),
		Folder_id:  tag.FolderID,
		Created_by: tag.CreatedBy,
	}

	var buf bytes.Buffer
	err := data.Serialize(&buf)
	if err != nil {
		return err
	}

	return s.publisher.Publish("Tag-Update", buf.Bytes())
}

func (s *TagService) DeleteTag(id string) error {
	data := map[string]interface{}{
		"id": id,
	}
	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}
	return s.publisher.Publish("Tag-Delete", encoded)
}

func (s *TagService) ListTags() ([]organization.Tag, error) {
	// TODO: Replace with real implementation
	return []organization.Tag{}, nil
}

func (s *TagService) GetTag(id string) (organization.Tag, error) {
	// TODO: Replace with real implementation
	return organization.Tag{Id: id}, nil
}
