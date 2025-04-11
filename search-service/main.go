package main

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/app"
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		config.HandleError(err)
	}
	app, err := app.NewApp(*conf)
	if err != nil {
		panic(err)
	}
	app.Start()
}
