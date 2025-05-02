package events
import "log"

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ICredentialEventController interface {
	Create(*kafka.Message) error
	Update(*kafka.Message) error
	Delete(*kafka.Message) error
}

type CredentialEventController struct {}

func NewCredentialEventController() ICredentialEventController {
	return &CredentialEventController{}
}

func (c *CredentialEventController) Create(message *kafka.Message) error {
	// TODO: Implement Create logic
	log.Printf("Received Create message: %v", message)
	return nil
}

func (c *CredentialEventController) Update(message *kafka.Message) error {
	// TODO: Implement Update logic
	log.Printf("Received Update message: %v", message)
	return nil
}

func (c *CredentialEventController) Delete(message *kafka.Message) error {
	// TODO: Implement Delete logic
	log.Printf("Received Delete message: %v", message)
	return nil
}
