package events

import (
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
	var createFolderMessage any
	err := c.srClient.GetValue(*message, &createFolderMessage)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements IFolderEventController.
func (c *FolderEventController) Delete(*kafka.Message) error {
	panic("unimplemented")
}

// Update implements IFolderEventController.
func (c *FolderEventController) Update(*kafka.Message) error {
	panic("unimplemented")
}
