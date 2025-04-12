package app

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/DO-2K23-26/polypass-microservices/search-service/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
)

type App struct {
	Config     config.Config
	GrpcServer *grpc.Server
}

func NewApp(Config config.Config) (*App, error) {
	_, err := infrastructure.NewElasticAdapter(Config.EsHost, Config.EsPassword)
	if err != nil {
		return nil, err
	}

	_, err = infrastructure.NewKafkaAdapter(Config.KafkaHost, Config.ClientId)
	if err != nil {
		return nil, err
	}

	GrpcServer, err := grpc.NewServer(nil, nil, nil, Config.Port)
	if err != nil {
		return nil, err
	}

	return &App{
		Config, GrpcServer,
	}, nil
}

func (app *App) Start() error {
	return app.GrpcServer.Start()
}
