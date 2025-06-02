package avroschemas

import (
	"github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/schemautils"
)

func NewEncoder(schemaRegistryURL, subject, schema string) (*schemautils.AvroEncoder, error) {
	return schemautils.NewAvroEncoder(schemaRegistryURL, subject, schema)
}
