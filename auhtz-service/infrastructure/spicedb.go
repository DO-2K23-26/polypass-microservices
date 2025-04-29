package infrastructure

import (
	"context"
	"log"
	"os"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
)

type SpiceDBAdapter struct {
	authzedClient *authzed.Client
}

func NewSpiceDBAdapter(authzedClient *authzed.Client) *SpiceDBAdapter {
	return &SpiceDBAdapter{
		authzedClient: authzedClient,
	}
}

func (s *SpiceDBAdapter) Close() error {
	return s.authzedClient.Close()
}

func (s *SpiceDBAdapter) HealthCheck() error {
	// Example health check: Attempt to read relationships
	ctx := context.Background()
	_, err := s.authzedClient.ReadSchema(ctx, &v1.ReadSchemaRequest{})
	if err != nil {
		log.Println("health check failed: %w", err)
		return err
	}
	return nil
}

func (s *SpiceDBAdapter) Init() error {
	schema, err := os.ReadFile("schema.zed")
	if err != nil {
		log.Printf("failed to read schema file: %v", err)
		return err
	}
	schemaString := string(schema)
	res , err := s.authzedClient.SchemaServiceClient.WriteSchema(context.Background(), &v1.WriteSchemaRequest{
		Schema: schemaString,
	})
	if err != nil {
		log.Println("Could not write the schema to spicedb:", err)
		return err
	}
	log.Println("Wrote the schema to spicedb", res)
	
	return nil
}
