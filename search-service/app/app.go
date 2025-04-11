package app

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/DO-2K23-26/polypass-microservices/search-service/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
)

type App struct {
	Config   config.Config
	EsClient *infrastructure.ElasticAdapter
	// UserRepository        userRepo.UserRepository
	// TagRepository         tagsRepo.TagRepository
	// FolderRepository      folderRepo.FolderRepository
	// CredentialsRepository credentialRepo.CredentialRepository
	// UserService           userService.UserService
	// TagService            tagsService.TagService
	// FolderService         folderService.FolderService
	// CredentialsService    credentialService.CredentialService
	GrpcServer *grpc.Server
}

func NewApp(Config config.Config) (*App, error) {
	EsClient, err := infrastructure.NewElasticAdapter(Config.EsHost, Config.EsPassword)
	if err != nil {
		return nil, err
	}

	GrpcServer, err := grpc.NewServer(nil, nil, nil, Config.Port)
	if err != nil {
		return nil, err
	}

	return &App{
		Config, EsClient, GrpcServer,
	}, nil
}


func (app *App) Start() error {
	return app.GrpcServer.Start()
}


