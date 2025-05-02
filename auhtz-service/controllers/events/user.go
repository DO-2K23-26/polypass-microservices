package events

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IUserEventController interface {
	AddUserToFolder(*kafka.Message) error
	RemoveUserToFolder(*kafka.Message) error
}

type UserEventController struct{}

func NewUserEventController() IUserEventController {
	return &UserEventController{}
}

func (c *UserEventController) AddUserToFolder(message *kafka.Message) error {
	// TODO: Implement AddUserToFolder logic
	log.Printf("Received AddUserToFolder message: %v", message)
	return nil
}

func (c *UserEventController) RemoveUserToFolder(message *kafka.Message) error {
	// TODO: Implement RemoveUserToFolder logic
	log.Printf("Received RemoveUserToFolder message: %v", message)
	return nil
}
