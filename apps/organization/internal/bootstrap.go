package internal

import (
	"fmt"
	"log"

	avroschemas "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas"
	generated "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"

	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/config"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/infrastructure/kafka"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/server"
	httpPorts "github.com/DO-2K23-26/polypass-microservices/organization/internal/transport/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	HttpServer *server.HttpServer
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	producer, err := kafka.NewProducerAdapter(cfg.KafkaHost, cfg.ClientId)
	if err != nil {
		return nil, err
	}

	// On utilise maintenant le Schema() de la struct générée
	folderEncoder, err := avroschemas.NewEncoder(cfg.SchemaRegistryURL, "organization-folder-event-value", generated.FolderEvent{}.Schema())
	if err != nil {
		return nil, err
	}

	tagEncoder, err := avroschemas.NewEncoder(cfg.SchemaRegistryURL, "organization-tag-event-value", generated.TagEvent{}.Schema())
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "localhost", "postgres", "postgres", "postgres", "5432")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("An error occured during the connection to the database :", err)
	}

	folderService := app.NewFolderService(producer, folderEncoder, db)
	tagService := app.NewTagService(producer, tagEncoder, db)

	folderHandler := httpPorts.NewFolderHandler(folderService)
	tagHandler := httpPorts.NewTagHandler(tagService)

	httpServer := server.NewHttpServer(cfg.HttpPort, folderHandler, tagHandler)

	return &App{
		HttpServer: httpServer,
	}, nil
}

func (a *App) Start() {
	a.HttpServer.Start()
}
