package main

import (
	"fmt"
	"os"

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
		fmt.Printf("Failed to create the application: %v\n", err)
		os.Exit(1)
	}
	err = app.Init()
	if err != nil {
		fmt.Printf("Failed to initialize the application: %v\n", err)
		os.Exit(1)
	}

	err = app.Start()
	if err != nil {
		fmt.Printf("Failed to start the application: %v\n", err)
		os.Exit(1)
	}
}
