package app

import (
	"log"

	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	httpController "github.com/DO-2K23-26/polypass-microservices/search-service/controller/http"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/http"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/kafka"
	credentialRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/credential"
	folderRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/folder"
	tagRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
	userRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
	credentialService"github.com/DO-2K23-26/polypass-microservices/search-service/services/credential"
	folderService "github.com/DO-2K23-26/polypass-microservices/search-service/services/folder"
	tagService "github.com/DO-2K23-26/polypass-microservices/search-service/services/tags"
	userService "github.com/DO-2K23-26/polypass-microservices/search-service/services/user"

	"github.com/DO-2K23-26/polypass-microservices/search-service/services/health"

	"sync"
)

type App struct {
	Config               config.Config
	esClient             *infrastructure.ElasticAdapter
	gormClient           *infrastructure.GormAdapter
	kafkaClient          *infrastructure.KafkaAdapter
	GrpcServer           *grpc.Server
	HttpServer           *http.Server
	KafkaConsumers       []kafka.KafkaConsumerConfig
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

	gormClient, err := infrastructure.NewGormAdapter(Config.PgHost, Config.PgUser, Config.PgPassword, Config.PgDBName, Config.PgPort)
	if err != nil {
		return nil, err
	}
	healthService := health.NewHealthService(esClient, kafkaClient, gormClient)
	healthController := httpController.NewHealthController(healthService)
	HttpServer := http.NewServer(healthController, Config.HttpPort)

	// Initialize repos
	userRepository := userRepository.NewGormUserRepository(gormClient.Db)
	folderRepository := folderRepository.NewEsSqlFolderRepository(gormClient, esClient)
	credentialRepository := credentialRepository.NewCredentialRepository(*esClient)
	tagRepository := tagRepository.NewESTagRepository(*esClient)

	// Initialize services
	credentialService := credentialService.NewCredentialService(credentialRepository)
	folderService := folderService.NewFolderService(folderRepository)
	tagService := tagService.NewTagService(tagRepository)
	userService := userService.NewUserService(userRepository)

	GrpcServer, err := grpc.NewServer(credentialService, folderService, tagService, Config.GrpcPort)
	if err != nil {
		return nil, err
	}

	// Create an instance of Consumers
	consumers := kafka.NewConsumers(
		credentialService,
		folderService,
		tagService,
		userService,
	)

	// Define Kafka consumers
	kafkaConsumers := []kafka.KafkaConsumerConfig{
		{Topic: "tag_creation", HandleMessage: consumers.HandleTagCreation, HandleError: kafka.HandleError},
		{Topic: "tag_deletion", HandleMessage: consumers.HandleTagDeletion, HandleError: kafka.HandleError},
		{Topic: "tag_update", HandleMessage: consumers.HandleTagUpdate, HandleError: kafka.HandleError},
		{Topic: "folder_creation", HandleMessage: consumers.HandleFolderCreation, HandleError: kafka.HandleError},
		{Topic: "folder_deletion", HandleMessage: consumers.HandleFolderDeletion, HandleError: kafka.HandleError},
		{Topic: "folder_update", HandleMessage: consumers.HandleFolderUpdate, HandleError: kafka.HandleError},
		{Topic: "credential_creation", HandleMessage: consumers.HandleCredentialCreation, HandleError: kafka.HandleError},
		{Topic: "credential_deletion", HandleMessage: consumers.HandleCredentialDeletion, HandleError: kafka.HandleError},
		{Topic: "credential_update", HandleMessage: consumers.HandleCredentialUpdate, HandleError: kafka.HandleError},
		{Topic: "user_creation", HandleMessage: consumers.HandleUserCreation, HandleError: kafka.HandleError},
		{Topic: "user_deletion", HandleMessage: consumers.HandleUserDeletion, HandleError: kafka.HandleError},
		{Topic: "user_update", HandleMessage: consumers.HandleUserUpdate, HandleError: kafka.HandleError},
	}

	return &App{
		Config:               Config,
		esClient:             esClient,
		gormClient:           gormClient,
		kafkaClient:          kafkaClient,
		GrpcServer:           GrpcServer,
		HttpServer:           HttpServer,
		KafkaConsumers:       kafkaConsumers,
	}, nil
}
func (app *App) Init() error {
	if err := app.esClient.CreateIndexes(); err != nil {
		log.Println("Could not create elastic indexes:", err)
		return err

	}
	if err := app.gormClient.Migrate(); err != nil {
		log.Println("Could not migrate database:", err)
		return err
	} else {
		log.Println("Postgres auto-migrated successfully")
	}

	return nil
}

func (app *App) Start() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(app.KafkaConsumers)+2)

	// Start Kafka consumers
	for _, consumer := range app.KafkaConsumers {
		wg.Add(1)
		go func(c kafka.KafkaConsumerConfig) {
			defer wg.Done()
			err := app.kafkaClient.Consume(c.Topic, c.HandleMessage, c.HandleError)
			if err != nil {
				errChan <- err
			}
		}(consumer)
	}

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
