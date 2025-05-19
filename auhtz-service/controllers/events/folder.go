package events

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/folder"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IFolderEventController interface {
	Create(*kafka.Message) error
	Delete(*kafka.Message) error
	Update(*kafka.Message) error
}

type FolderEventController struct {
	folderService *folder.FolderService
}

func NewFolderEventController(folderService *folder.FolderService) IFolderEventController {
	return &FolderEventController{
		folderService: folderService,
	}
}

func (c *FolderEventController) Create(message *kafka.Message) error {
	log.Printf("Received Create message: %v", message)

	// Parse the message to extract folderId and parentId
	folderId := "extracted_folder_id" // Replace with actual parsing logic
	parentId := "extracted_parent_id" // Replace with actual parsing logic

	return c.folderService.Create(context.Background(), folderId, parentId)
}

func (c *FolderEventController) Delete(message *kafka.Message) error {
	log.Printf("Received Delete message: %v", message)

	// Parse the message to extract folderId
	folderId := "extracted_folder_id" // Replace with actual parsing logic

	return c.folderService.Delete(context.Background(), folderId)
}

func (c *FolderEventController) Update(message *kafka.Message) error {
	log.Printf("Received Update message: %v", message)

	// Parse the message to extract folderId and newParentId
	folderId := "extracted_folder_id"       // Replace with actual parsing logic
	newParentId := "extracted_new_parent_id" // Replace with actual parsing logic

	return c.folderService.UpdateParent(context.Background(), folderId, newParentId)
}
