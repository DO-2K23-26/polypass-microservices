package app

import (
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
	data := map[string]interface{}{
		"id":         tag.Id,
		"name":       tag.Name,
		"color":      tag.Color,
		"created_at": tag.CreatedAt,
		"updated_at": tag.UpdatedAt,
		"folder_id":  tag.FolderId,
		"created_by": tag.CreatedBy,
	}

	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}

	return s.publisher.Publish("Tag-Create", encoded)
}

func (s *TagService) UpdateTag(tag organization.Tag) error {
	data := map[string]interface{}{
		"id":         tag.Id,
		"name":       tag.Name,
		"color":      tag.Color,
		"created_at": tag.CreatedAt,
		"updated_at": tag.UpdatedAt,
		"folder_id":  tag.FolderId,
		"created_by": tag.CreatedBy,
	}
	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}
	return s.publisher.Publish("Tag-Update", encoded)
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
