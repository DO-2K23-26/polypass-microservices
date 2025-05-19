package app

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/config"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/controllers/events"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/internal/consumers"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/producer"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/credential"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/folder"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/tag"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/services/user"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	authzedClient *infrastructure.SpiceDBAdapter
	consumers     *consumers.Consumers
	producer      *producer.Producer
}

func NewApp(config config.Config) (*App, error) {
	authzedRawClient, err := authzed.NewClient(config.AuthzedHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(config.AuthzedApiKey))
	if err != nil {
		return nil, err
	}
	kafkaClient, err := infrastructure.NewKafkaAdapter(config.KafkaHost, config.KafkaClientId)
	if err != nil {
		return nil, err
	}

	// Instantiate producer for debug purpose.
	// producer := producer.NewProducer(*kafkaClient)

	authzedClient := infrastructure.NewSpiceDBAdapter(authzedRawClient)
	// Instantiate services
	folderService := folder.NewFolderService(authzedClient)
	credentialService := credential.NewCredentialService(authzedClient)
	userService := user.NewUserService(authzedClient)
	tagService := tag.NewTagService(authzedClient)
	
	
	
	// Instantiate the event controllers
	folderEventController := events.NewFolderEventController(folderService)
	credentialEventController := events.NewCredentialEventController(credentialService)
	userEventController := events.NewUserEventController(userService)
	tagEventController := events.NewTagEventController(tagService)

	consumersController := consumers.NewConsumers(folderEventController, credentialEventController, tagEventController, userEventController, *kafkaClient)
	return &App{
		// producer:      producer,
		consumers:     consumersController,
		authzedClient: authzedClient,
	}, nil
}

// Perform instanciation to external services/ local services/ repos
func (a *App) Start() error {
	ctx := context.Background()
	err := a.consumers.Start(ctx)
	if err != nil {
		log.Fatal("Error while starting:", err)
		return err
	} else {
		log.Println("Consumers started")
	}
	<-ctx.Done()
	return nil
}

// Perform any data migration the version of the app need
func (a *App) Init() error {
	if err := a.authzedClient.Migrate(); err != nil {
		log.Println("Error while migrating:", err)
		return err
	}
	return nil
}

// Allow to stop the app gracefully
func (a *App) Stop() error {
	a.authzedClient.Close()
	return nil
}
