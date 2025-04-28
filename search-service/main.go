package main

import (
	"fmt"
	"os"

	"github.com/DO-2K23-26/polypass-microservices/search-service/app"
	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/user"
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
	_ , err = app.UserRepository.Create(user.CreateUserQuery{User: types.User{
		ID:      "test",
		Folders: []types.Folder{},
	}})
	
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		os.Exit(1)
	}
	
	_, err = app.UserRepository.Get(user.GetUserQuery{ID: "test"})
	
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		os.Exit(1)
	}
	
	err = app.Start()
	if err != nil {
		fmt.Printf("Failed to start the application: %v\n", err)
		os.Exit(1)
	}
}
