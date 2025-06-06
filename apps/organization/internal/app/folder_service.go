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

	return s.publisher.Publish("Folder-Create", encoded)
}

func (s *FolderService) UpdateFolder(folder organization.Folder) error {
	data := map[string]interface{}{
		"id":   folder.Id,
		"name": folder.Name,
	}
	encoded, err := s.encoder.Encode(data)
	if err != nil {
		return err
	}
	return s.publisher.Publish("Folder-Update", encoded)
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
