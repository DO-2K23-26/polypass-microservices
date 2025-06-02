package app

import (
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
	data := map[string]interface{}{
		"id":   folder.Id,
		"name": folder.Name,
	}

	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}

	return s.publisher.Publish("create_folder", encoded)
}
