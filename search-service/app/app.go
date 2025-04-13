package app

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
)

type App struct {
	Config     config.Config
	esClient   *infrastructure.ElasticAdapter
	GrpcServer *grpc.Server
}

func NewApp(Config config.Config) (*App, error) {
	esClient, err := infrastructure.NewElasticAdapter(Config.EsHost, Config.EsUsername, Config.EsPassword)
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
		Config, esClient, GrpcServer,
	}, nil
}

func (app *App) Init() error {
	app.esClient.CreateIndexes()
	return nil
}

func (app *App) Start() error {
	return app.GrpcServer.Start()
}
