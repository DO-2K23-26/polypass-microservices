package events

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ITagEventController interface {
	Create(*kafka.Message) error
	Update(*kafka.Message) error
	Delete(*kafka.Message) error
}

type TagEventController struct{}

func NewTagEventController() ITagEventController {
	return &TagEventController{}
}

func (c *TagEventController) Create(message *kafka.Message) error {
	// TODO: Implement Create logic
	log.Printf("Received Create message: %v", message)
	return nil
}

func (c *TagEventController) Update(message *kafka.Message) error {
	// TODO: Implement Update logic
	log.Printf("Received Update message: %v", message)
	return nil
}

func (c *TagEventController) Delete(message *kafka.Message) error {
	// TODO: Implement Delete logic
	log.Printf("Received Delete message: %v", message)
	return nil
}
