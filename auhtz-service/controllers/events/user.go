package events

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/user"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IUserEventController interface {
	AddUserToFolder(*kafka.Message) error
	RemoveUserToFolder(*kafka.Message) error
}

type UserEventController struct {
	UserService *user.UserService
}

func NewUserEventController(userService *user.UserService) IUserEventController {
	return &UserEventController{
		UserService: userService,
	}
}

func (c *UserEventController) AddUserToFolder(message *kafka.Message) error {
	// TODO: Implement AddUserToFolder logic
	log.Printf("Received AddUserToFolder message: %v", message)
	userId := "extracted_user_id"     // Replace with actual parsing logic
	folderId := "extracted_folder_id" // Replace with actual parsing logic
	relation := "extracted_relation"  // Replace with actual parsing logic
	err := c.UserService.AddUserToFolder(context.Background(), userId, folderId, relation)
	if err != nil {
		log.Printf("Failed to add user to folder: %v", err)
		return err
	}
	return nil
}

func (c *UserEventController) RemoveUserToFolder(message *kafka.Message) error {
	// TODO: Implement RemoveUserToFolder logic
	log.Printf("Received RemoveUserToFolder message: %v", message)
	userId := "extracted_user_id"     // Replace with actual parsing logic
	folderId := "extracted_folder_id" // Replace with actual parsing logic
	err := c.UserService.RemoveUserFromFolder(context.Background(), userId, folderId)
	if err != nil {
		log.Printf("Failed to remove user from folder: %v", err)
		return err
	}
	return nil
}
