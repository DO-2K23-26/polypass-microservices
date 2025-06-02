package internal

import (
	avroschemas "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas"
	generated "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/generated"

	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/config"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/infrastructure/kafka"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/server"
	httpPorts "github.com/DO-2K23-26/polypass-microservices/organization/internal/transport/http"
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

	folderService := app.NewFolderService(producer, folderEncoder)
	tagService := app.NewTagService(producer, tagEncoder)

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
