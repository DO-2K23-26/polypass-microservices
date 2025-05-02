package consumers

import (
	"context"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/controllers/events"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type Consumers struct {
	kafka                     infrastructure.KafkaAdapter
	FolderEventController     events.IFolderEventController
	CredentialEventController events.ICredentialEventController
	TagEventController        events.ITagEventController
	UserEventController       events.IUserEventController
}

func NewConsumers(folderEventController events.IFolderEventController, credentialEventController events.ICredentialEventController, tagEventController events.ITagEventController, userEventController events.IUserEventController, kafka infrastructure.KafkaAdapter) *Consumers {
	return &Consumers{
		kafka:                     kafka,
		FolderEventController:     folderEventController,
		CredentialEventController: credentialEventController,
		TagEventController:        tagEventController,
		UserEventController:       userEventController,
	}
}

func (c *Consumers) Start(ctx context.Context) error {
	// Register consumers
	// Folder event
	err := c.kafka.Consume("create_folder", c.FolderEventController.Create, nil, ctx)
	err = c.kafka.Consume("delete_folder", c.FolderEventController.Delete, nil, ctx)
	err = c.kafka.Consume("update_folder", c.FolderEventController.Update, nil, ctx)

	// Credentials event
	err = c.kafka.Consume("create_credential", c.CredentialEventController.Create, nil, ctx)
	err = c.kafka.Consume("update_credential", c.CredentialEventController.Update, nil, ctx)
	err = c.kafka.Consume("delete_credential", c.CredentialEventController.Delete, nil, ctx)

	// Tags event
	err = c.kafka.Consume("create_tag", c.TagEventController.Create, nil, ctx)
	err = c.kafka.Consume("update_tag", c.TagEventController.Update, nil, ctx)
	err = c.kafka.Consume("delete_tag", c.TagEventController.Delete, nil, ctx)

	// Users event
	err = c.kafka.Consume("add_user", c.UserEventController.AddUserToFolder, nil, ctx)
	err = c.kafka.Consume("revoke_user", c.UserEventController.AddUserToFolder, nil, ctx)
	return err
}
