package events

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/tag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ITagEventController interface {
	Create(*kafka.Message) error
	Update(*kafka.Message) error
	Delete(*kafka.Message) error
}

type TagEventController struct {
	tagService *tag.TagService
}

func NewTagEventController(tagService *tag.TagService) ITagEventController {
	return &TagEventController{
		tagService: tagService,
	}
}

func (c *TagEventController) Create(message *kafka.Message) error {
	log.Printf("Received Create message: %v", message)
	tagId := "extracted_tag_id"     // Replace with actual parsing logic
	folderID := "extracted_folder_id" // Replace with actual parsing logic
	err := c.tagService.Create(context.Background(), tagId, folderID)
	if err != nil {
		log.Printf("Error handling Create message: %v", err)
		return err
	}
	return nil
}

func (c *TagEventController) Update(message *kafka.Message) error {
	log.Printf("Received Update message: %v", message)
	tagId := "extracted_tag_id"     // Replace with actual parsing logic
	folderID := "extracted_folder_id" // Replace with actual parsing logic
	err := c.tagService.Update(context.Background(), tagId, folderID)
	if err != nil {
		log.Printf("Error handling Update message: %v", err)
		return err
	}
	return nil
}

func (c *TagEventController) Delete(message *kafka.Message) error {
	log.Printf("Received Delete message: %v", message)
	tagId := "extracted_tag_id" // Replace with actual parsing logic
	err := c.tagService.Delete(context.Background(), tagId)
	if err != nil {
		log.Printf("Error handling Delete message: %v", err)
		return err
	}
	return nil
}
