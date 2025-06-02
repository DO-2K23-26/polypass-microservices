package avro

import (
    "github.com/riferrei/srclient"
    "github.com/DO-2K23-26/polypass-microservices/libs/schemautils"
)

func NewEncoder(schemaRegistryURL, subject string, schema string) (*schemautils.AvroEncoder, error) {
    client := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
    return schemautils.NewAvroEncoder(client, subject, schema)
}
