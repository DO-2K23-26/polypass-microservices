package consumers

import (
	"github.com/DO-2K23-26/polypass-microservices/authz-service/controllers/events"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type Consumers struct {
	kafka                     infrastructure.KafkaAdapter
	srclient                  infrastructure.SchemaRegistry
	FolderEventController     events.IFolderEventController
	CredentialEventController events.ICredentialEventController
	TagEventController        events.ITagEventController
	UserEventController       events.IUserEventController
}

func NewConsumers(folderEventController events.IFolderEventController) *Consumers {
	return &Consumers{
		FolderEventController: folderEventController,
	}
}

func (c *Consumers) Start() error {
	// Register consumers
	// Folder event
	c.kafka.Consume("create_folder", c.FolderEventController.Create, nil)
	c.kafka.Consume("delete_folder", c.FolderEventController.Delete, nil)
	c.kafka.Consume("update_folder", c.FolderEventController.Update, nil)
	
	// Credentials event
	c.kafka.Consume("create_credential", c.CredentialEventController.Create, nil)
	c.kafka.Consume("update_credential", c.CredentialEventController.Update, nil)
	c.kafka.Consume("delete_credential", c.CredentialEventController.Delete, nil)
	
	// Tags event
	c.kafka.Consume("create_tag", c.TagEventController.Create, nil)
	c.kafka.Consume("update_tag", c.TagEventController.Update, nil)
	c.kafka.Consume("delete_tag", c.TagEventController.Delete, nil)
	
	// Users event
	c.kafka.Consume("add_user", c.UserEventController.AddUserToFolder, nil)
	c.kafka.Consume("revoke_user", c.UserEventController.AddUserToFolder, nil)
	return nil
}
