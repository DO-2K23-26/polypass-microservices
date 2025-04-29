package app

import (
	"github.com/DO-2K23-26/polypass-microservices/authz-service/config"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	authzedClient *infrastructure.SpiceDBAdapter
}

func NewApp(config config.Config) (*App, error) {
	authzedRawClient, err := authzed.NewClient(config.AuthzedHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(config.AuthzedApiKey))
	if err != nil {
		return nil, err
	}
	
	authzedClient:= infrastructure.NewSpiceDBAdapter(authzedRawClient)
	
	return &App{
		authzedClient: authzedClient,
	}, nil
}

// Perform instanciation to external services/ local services/ repos
func (a *App) Start() error {

	return nil
}

// Perform any data migration the version of the app need
func (a *App) Init() error {
	// ...
	return nil
}

// Allow to stop the app gracefully
func (a *App) Stop() error {
	a.authzedClient.Close()
	return nil
}
