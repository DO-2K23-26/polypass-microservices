package app

import (
	"log"

	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	httpController "github.com/DO-2K23-26/polypass-microservices/search-service/controller/http"
	"github.com/DO-2K23-26/polypass-microservices/search-service/infrastructure"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/grpc"
	"github.com/DO-2K23-26/polypass-microservices/search-service/internal/api/http"
	credentialRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/credential"
	folderRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/folder"
	tagRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
	userRepository "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"

	"github.com/DO-2K23-26/polypass-microservices/search-service/services/health"

	"sync"
)

type App struct {
	Config               config.Config
	esClient             *infrastructure.ElasticAdapter
	gormClient           *infrastructure.GormAdapter
	UserRepository       *userRepository.GormUserRepository
	FolderRepository     *folderRepository.IFolderRepository
	CredentialRepository *credentialRepository.ICredentialRepository
	TagRepository        *tagRepository.ITagRepository
	GrpcServer           *grpc.Server
	HttpServer           *http.Server
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

	GrpcServer, err := grpc.NewServer(nil, nil, nil, Config.GrpcPort)
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
	return &App{
		Config:               Config,
		esClient:             esClient,
		gormClient:           gormClient,
		GrpcServer:           GrpcServer,
		HttpServer:           HttpServer,
		UserRepository:       userRepository,
		FolderRepository:     &folderRepository,
		TagRepository:        &tagRepository,
		CredentialRepository: &credentialRepository,
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
