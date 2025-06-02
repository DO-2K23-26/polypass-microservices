package app

import (
    "github.com/DO-2K23-26/polypass-microservices/interfaces/organization"
    "github.com/DO-2K23-26/polypass-microservices/avro-schemas/schemautils"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/infrastructure"
)

type FolderService struct {
    kafka   *infrastructure.KafkaAdapter
    encoder *schemautils.AvroEncoder
}

func NewFolderService(kafka *infrastructure.KafkaAdapter, encoder *schemautils.AvroEncoder) *FolderService {
    return &FolderService{
        kafka:   kafka,
        encoder: encoder,
    }
}

func (s *FolderService) CreateFolder(folder organization.Folder) error {
    // Convertir le modèle en map pour Avro
    data := map[string]interface{}{
        "id":   folder.Id,
        "name": folder.Name,
        // si besoin, ajouter d'autres champs compatibles avec ton schéma Avro
    }
    return s.kafka.ProduceAvro("create_folder", s.encoder, data)
}
