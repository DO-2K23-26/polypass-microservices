package main

import (
	"os"

	"github.com/DO-2K23-26/polypass-microservices/credentials/application/http"
	"github.com/DO-2K23-26/polypass-microservices/credentials/config"
	"github.com/DO-2K23-26/polypass-microservices/credentials/core"
	"github.com/DO-2K23-26/polypass-microservices/credentials/infrastructure/sql"
	"github.com/optique-dev/optique"
)

// @title Polypass Credentials Microservice
// @version 0.1.0
// @description Polypass Credentials Microservice
// @contact.name Tristan-Mihai Radulescu
// @contact.url https://github.com/DO-2K23-26
// @contact.email tristan-mihai.radulescu@etu.umontpellier.fr
func main() {
	conf, err := config.LoadConfig()

	if err != nil {
		config.HandleError(err)
	}
	cycle := NewCycle()

	database, err := sql.NewSql(conf.Database)

	if err != nil {
		optique.Error(err.Error())
		cycle.Stop()
		os.Exit(1)
	}


	// service
	credential_service := core.NewCredentialService(database)

	// controllers
	credentials_controller := http.NewCredentialsController(credential_service)
	docs_controller := http.NewDocsController()
	health_controller := http.NewHealthController()

	http_server, err := http.NewHttp(conf.Server)

	if err != nil {
		optique.Error(err.Error())
		cycle.Stop()
		os.Exit(1)
	}
	http_server.WithHandler(credentials_controller)
	http_server.WithHandler(docs_controller)
	http_server.WithHandler(health_controller)

	cycle.AddApplication(http_server)

	cycle.AddRepository(database)

	if conf.Bootstrap {
		err := cycle.Setup()
		if err != nil {
			optique.Error(err.Error())
			cycle.Stop()
			os.Exit(1)
		}
	}

	err = cycle.Ignite()
}
