package kafka

import (
	"log"
	"os"
	"path/filepath"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/hamba/avro"

	credentialService "github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	folderService "github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tagService "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
)

// ConsumerInterface defines the interface for Kafka consumers
type IConsumers interface {
	HandleTagCreation(msg *kafka.Message) error
	HandleTagDeletion(msg *kafka.Message) error
	HandleTagUpdate(msg *kafka.Message) error
	HandleFolderCreation(msg *kafka.Message) error
	HandleFolderDeletion(msg *kafka.Message) error
	HandleFolderUpdate(msg *kafka.Message) error
	HandleCredentialCreation(msg *kafka.Message) error
	HandleCredentialDeletion(msg *kafka.Message) error
	HandleCredentialUpdate(msg *kafka.Message) error
}

// Ensure Consumers implements ConsumerInterface
var _ IConsumers = &Consumers{} // This is a common go hack to make the compiler assert that Consumers implements the IConsumers interface, without needing to call it explicitly anywhere in the code.

// SearchServiceServer implements the gRPC search service
type Consumers struct {
	credentialService *credentialService.CredentialService
	folderService     *folderService.FolderService
	tagService        *tagService.TagService
}

// NewConsumers contains the set of services used by Kafka messages handlers
func NewConsumers(
	credentialService *credentialService.CredentialService,
	folderService *folderService.FolderService,
	tagService *tagService.TagService,
) *Consumers {
	if credentialService == nil || folderService == nil || tagService == nil {
		log.Fatal("Services must not be nil")
	}
	return &Consumers{
		credentialService: credentialService,
		folderService:     folderService,
		tagService:        tagService,
	}
}

// KafkaConsumerConfig defines the configuration for a Kafka consumer
type KafkaConsumerConfig struct {
	Topic         string
	HandleMessage func(*kafka.Message) error
	HandleError   func(error)
}

// LoadAvroSchema loads an Avro schema from a file
func LoadAvroSchema(schemaPath string) (avro.Schema, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}
	schema, err := avro.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// Load schemas from the `interfaces/avro` folder (Join is used to ensure the path is correct regardless of OS)
var (
	tagCreationSchema, _        = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "tag_event_created.avsc"))
	tagDeletionSchema, _        = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "tag_event_deleted.avsc"))
	tagUpdateSchema, _          = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "tag_event_updated.avsc"))
	folderCreationSchema, _     = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "folder_event_created.avsc"))
	folderDeletionSchema, _     = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "folder_event_deleted.avsc"))
	folderUpdateSchema, _       = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "folder_event_updated.avsc"))
	credentialCreationSchema, _ = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "credential_event_created.avsc"))
	credentialDeletionSchema, _ = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "credential_event_deleted.avsc"))
	credentialUpdateSchema, _   = LoadAvroSchema(filepath.Join("..", "..", "..", "interfaces", "avro", "credential_event_updated.avsc"))
)

// Tag Handlers
func (c Consumers) HandleTagCreation(msg *kafka.Message) error {
	log.Println("Handling tag creation:", string(msg.Value))

	var req tagService.CreateTagRequest
	if err := avro.Unmarshal(tagCreationSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing tag creation message:", err)
		return err
	}

	_, err := c.tagService.Create(req)
	return err
}

func (c Consumers) HandleTagDeletion(msg *kafka.Message) error {
	log.Println("Handling tag deletion:", string(msg.Value))

	var req tagService.DeleteTagRequest
	if err := avro.Unmarshal(tagDeletionSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing tag deletion message:", err)
		return err
	}

	err := c.tagService.Delete(req)
	return err
}

func (c Consumers) HandleTagUpdate(msg *kafka.Message) error {
	log.Println("Handling tag update:", string(msg.Value))

	var req tagService.UpdateTagRequest
	if err := avro.Unmarshal(tagUpdateSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing tag update message:", err)
		return err
	}

	err := c.tagService.Update(req)
	return err
}

// Folder Handlers
func (c Consumers) HandleFolderCreation(msg *kafka.Message) error {
	log.Println("Handling folder creation:", string(msg.Value))

	var req folderService.CreateFolderRequest
	if err := avro.Unmarshal(folderCreationSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing folder creation message:", err)
		return err
	}

	_, err := c.folderService.Create(req)
	return err
}

func (c Consumers) HandleFolderDeletion(msg *kafka.Message) error {
	log.Println("Handling folder deletion:", string(msg.Value))

	var req folderService.DeleteFolderRequest
	if err := avro.Unmarshal(folderDeletionSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing folder deletion message:", err)
		return err
	}

	err := c.folderService.Delete(req)
	return err
}

func (c Consumers) HandleFolderUpdate(msg *kafka.Message) error {
	log.Println("Handling folder update:", string(msg.Value))

	var req folderService.UpdateFolderRequest
	if err := avro.Unmarshal(folderUpdateSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing folder update message:", err)
		return err
	}

	_, err := c.folderService.Update(req)
	return err
}

// Credential Handlers
func (c Consumers) HandleCredentialCreation(msg *kafka.Message) error {
	log.Println("Handling credential creation:", string(msg.Value))

	var req credentialService.CreateCredentialRequest
	if err := avro.Unmarshal(credentialCreationSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing credential creation message:", err)
		return err
	}

	_, err := c.credentialService.Create(req)
	return err
}

func (c Consumers) HandleCredentialDeletion(msg *kafka.Message) error {
	log.Println("Handling credential deletion:", string(msg.Value))

	var req credentialService.DeleteCredentialRequest
	if err := avro.Unmarshal(credentialDeletionSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing credential deletion message:", err)
		return err
	}

	err := c.credentialService.Delete(req)
	return err
}

func (c Consumers) HandleCredentialUpdate(msg *kafka.Message) error {
	log.Println("Handling credential update:", string(msg.Value))

	var req credentialService.UpdateCredentialRequest
	if err := avro.Unmarshal(credentialUpdateSchema, msg.Value, &req); err != nil {
		log.Println("Error deserializing credential update message:", err)
		return err
	}

	err := c.credentialService.Update(req)
	return err
}

// Error Handler
func HandleError(err error) {
	log.Println("Error handling Kafka message:", err)
	// Could implement retry logic or other error handling strategies here,
	// such as sending the error to a dead-letter queue
	// or structured logging for analysis.
}
