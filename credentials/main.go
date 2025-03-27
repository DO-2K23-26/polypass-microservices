package main

import (
	"log/slog"
	"os"

	"github.com/DO-2K23-26/polypass-microservices/credentials/config"
  "github.com/DO-2K23-26/polypass-microservices/credentials/cycle"
)

// @title Optique application TO CHANGE
// @version 1.0
// @description This is a sample application
// @contact.name Courtcircuits
// @contact.url https://github.com/Courtcircuits
// @contact.email tristan-mihai.radulescu@etu.umontpellier.fr
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		config.HandleError(err)
	}
	app_cycle := cycle.NewCycle()

	if conf.Bootstrap {
		err := app_cycle.Setup()
		if err != nil {
			slog.Error(err.Error())
			app_cycle.Stop()
			os.Exit(1)
		}
	}

	err = app_cycle.Ignite()
}
