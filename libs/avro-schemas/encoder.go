package avroschemas

import (
    "github.com/riferrei/srclient"
    "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
)

func NewEncoder(schemaRegistryURL, subject, schema string) (*schemautils.AvroEncoder, error) {
    client := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
    return schemautils.NewAvroEncoder(client, subject, schema)
}
