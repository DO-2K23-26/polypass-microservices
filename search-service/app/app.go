package app

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	httpController "github.com/DO-2K23-26/polypass-microservices/search-service/controller/http"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/http"
	"github.com/DO-2K23-26/polypass-microservices/search-service/services/health"

	"sync"
)

type App struct {
	Config     config.Config
	esClient   *infrastructure.ElasticAdapter
	GrpcServer *grpc.Server
	HttpServer *http.Server
}

func NewApp(Config config.Config) (*App, error) {
	esClient, err := infrastructure.NewElasticAdapter(Config.EsHost, Config.EsUsername, Config.EsPassword)
	if err != nil {
		return nil, err
	}

	kafkaClient, err := infrastructure.NewKafkaAdapter(Config.KafkaHost, Config.ClientId)
	if err != nil {
		return nil, err
	}

	GrpcServer, err := grpc.NewServer(nil, nil, nil, Config.GrpcPort)
	if err != nil {
		return nil, err
	}
	healthService := health.NewHealthService(esClient, kafkaClient)
	healthController := httpController.NewHealthController(healthService)
	HttpServer := http.NewServer(healthController, Config.HttpPort)

	return &App{
		Config, esClient, GrpcServer, HttpServer,
	}, nil
}

func (app *App) Init() error {
	app.esClient.CreateIndexes()
	return nil
}

func (app *App) Start() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.GrpcServer.Start(); err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.HttpServer.Start(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}
