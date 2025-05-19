package events

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/credential"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ICredentialEventController interface {
	Create(*kafka.Message) error
	Update(*kafka.Message) error
	Delete(*kafka.Message) error
}

type CredentialEventController struct {
	CredentialService *credential.CredentialService
}

func NewCredentialEventController(credentialService *credential.CredentialService) ICredentialEventController {
	return &CredentialEventController{
		CredentialService: credentialService,
	}
}

func (c *CredentialEventController) Create(message *kafka.Message) error {
	log.Printf("Received Create message: %v", message)

	// Parse the message to extract credential details
	credentialID := "extracted_credential_id" // Replace with actual parsing logic
	folderID := "extracted_folder_id"         // Replace with actual parsing logic

	// Pass parsed data to the service
	return c.CredentialService.Create(context.Background(), credentialID, folderID)
}

func (c *CredentialEventController) Update(message *kafka.Message) error {
	log.Printf("Received Update message: %v", message)

	// Parse the message to extract credential details
	credentialID := "extracted_credential_id" // Replace with actual parsing logic
	folderID := "extracted_folder_id"         // Replace with actual parsing logic

	// Pass parsed data to the service
	return c.CredentialService.Update(context.Background(), credentialID, folderID)
}

func (c *CredentialEventController) Delete(message *kafka.Message) error {
	log.Printf("Received Delete message: %v", message)

	// Parse the message to extract credential details
	credentialID := "extracted_credential_id" // Replace with actual parsing logic

	// Pass parsed data to the service
	return c.CredentialService.Delete(context.Background(), credentialID)
}
