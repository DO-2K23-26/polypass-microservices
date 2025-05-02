package events

import (
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IFolderEventController interface {
	Create(*kafka.Message) error
	Delete(*kafka.Message) error
	Update(*kafka.Message) error
}

type FolderEventController struct {
	srClient infrastructure.SchemaRegistry
}

func NewFolderEventController() IFolderEventController {
	return &FolderEventController{}
}

func (c *FolderEventController) Create(message *kafka.Message) error {
	// TODO: Implement Create logic
	log.Printf("Received Create message: %v", message)
	return nil
}

func (c *FolderEventController) Delete(message *kafka.Message) error {
	// TODO: Implement Delete logic
	log.Printf("Received Delete message: %v", message)
	return nil
}

func (c *FolderEventController) Update(message *kafka.Message) error {
	// TODO: Implement Update logic
	log.Printf("Received Update message: %v", message)
	return nil
}
