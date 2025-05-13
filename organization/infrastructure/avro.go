package infrastructure

import (
	"github.com/DO-2K23-26/polypass-microservices/avro-schemas/folder"
	"github.com/DO-2K23-26/polypass-microservices/avro-schemas/schemautils"
)

func NewOrganizationEncoder(schemaRegistryURL string) (*schemautils.AvroEncoder, error) {
	subject := "organization-topic-value"
	return schemautils.NewAvroEncoder(schemaRegistryURL, subject, folder.Schema)
}
