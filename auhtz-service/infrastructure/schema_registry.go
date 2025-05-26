package infrastructure

import (
	"encoding/binary"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/riferrei/srclient"
)

type SchemaRegistry struct {
	client *srclient.SchemaRegistryClient
}

func NewSchemaRegistry(url string) (*SchemaRegistry, error) {
	client := srclient.NewSchemaRegistryClient(url)

	return &SchemaRegistry{client: client}, nil
}

func (sr *SchemaRegistry) GetValue(message kafka.Message, result any) error {
	schemaID := binary.BigEndian.Uint32(message.Value[1:5])
	schema, err := sr.client.GetSchema(int(schemaID))
	if err != nil {
		return err
	}
	native, _, err := schema.Codec().NativeFromBinary(message.Value[5:])
	if err != nil {
		return err
	}

	// Attempt to cast the native value to the desired type
	if castedResult, ok := native.(any); ok {
		*result.(*any) = castedResult
		return nil
	}
	return fmt.Errorf("failed to cast value to the desired type")
}

func (sr *SchemaRegistry) CheckHealth() error {
	_, err := sr.client.GetSubjects()
	if err != nil {
		return fmt.Errorf("schema registry health check failed: %w", err)
	}
	return nil
}
