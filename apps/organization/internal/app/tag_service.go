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
		"id":   tag.Id,
		"name": tag.Name,
	}

	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}

	return s.publisher.Publish("create_tag", encoded)
}
