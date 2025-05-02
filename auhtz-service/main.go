package main

import (
	"log"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/app"
	"github.com/DO-2K23-26/polypass-microservices/authz-service/config"
)

func main() {
	config := config.LoadConfig()
	app, err := app.NewApp(*config)
	if err != nil {
		log.Fatal(err)
	}
	err = app.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = app.Start()
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	app.Stop()
}
