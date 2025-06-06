package main

import (
	"os"

	"github.com/DO-2K23-26/polypass-microservices/credentials/config"
	"github.com/DO-2K23-26/polypass-microservices/credentials/infrastructure/sql"
	"github.com/optique-dev/core"
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
		core.Error(err.Error())
		cycle.Stop()
		os.Exit(1)
	}

	cycle.AddRepository(database)

	if conf.Bootstrap {
		err := cycle.Setup()
		if err != nil {
			core.Error(err.Error())
			cycle.Stop()
			os.Exit(1)
		}
	}

	err = cycle.Ignite()
}
