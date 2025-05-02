package app

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/config"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/controllers/events"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/internal/consumers"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	authzedClient *infrastructure.SpiceDBAdapter
	consumers     *consumers.Consumers
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

	authzedClient := infrastructure.NewSpiceDBAdapter(authzedRawClient)
	folderEventController := events.NewFolderEventController()
	credentialEventController := events.NewCredentialEventController()
	tagEventController := events.NewTagEventController()
	userEventController := events.NewUserEventController()

	consumersController := consumers.NewConsumers(folderEventController, credentialEventController, tagEventController, userEventController, *kafkaClient)
	return &App{
		consumers:     consumersController,
		authzedClient: authzedClient,
	}, nil
}

// Perform instanciation to external services/ local services/ repos
func (a *App) Start() error {
	ctx := context.Background()
	a.consumers.Start(ctx)
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
