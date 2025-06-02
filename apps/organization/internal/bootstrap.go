package internal

import (
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/config"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/infrastructure/avro"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/infrastructure/kafka"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/server"
    httpPorts "github.com/DO-2K23-26/polypass-microservices/organization/internal/ports/http"

    "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/folder"
    "github.com/DO-2K23-26/polypass-microservices/libs/avro-schemas/tag"
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

    folderEncoder, err := avro.NewEncoder(cfg.SchemaRegistryURL, "create_folder-value", folder.Schema)
    if err != nil {
        return nil, err
    }

    tagEncoder, err := avro.NewEncoder(cfg.SchemaRegistryURL, "create_tag-value", tag.Schema)
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
