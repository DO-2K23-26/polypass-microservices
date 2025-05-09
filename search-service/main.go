package main

import (
	"fmt"
	"os"

	"github.com/DO-2K23-26/polypass-microservices/search-service/app"
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		config.HandleError(err)
	}
	application, err := app.NewApp(*conf)
	if err != nil {
		fmt.Printf("Failed to create the application: %v\n", err)
		os.Exit(1)
	}
	err = application.Init()
	if err != nil {
		fmt.Printf("Failed to initialize the application: %v\n", err)
		os.Exit(1)
	}
	
	tagRepo := *application.TagRepository

	// // Create a new tag
	// createQuery := tags.CreateTagQuery{
	// 	ID:       "test-tag-id",
	// 	Name:     "Test Tag",
	// 	FolderID: "test-folder-id",
	// }
	// createResult, err := tagRepo.Create(createQuery)
	// if err != nil {
	// 	fmt.Printf("Failed to create tag: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Tag created: %+v\n", createResult.Tag)

	// Use mGet to retrieve the created tag
	mGetQuery := tags.MGetTagQuery{
		IDs: []string{"test-tag-id"},
	}
	mGetResult, err := tagRepo.MGet(mGetQuery)
	if err != nil {
		fmt.Printf("Failed to retrieve tag with mGet: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Tags retrieved with mGet: %+v\n", mGetResult.Tags)
	// err = application.Start()
	// if err != nil {
	// 	fmt.Printf("Failed to start the application: %v\n", err)
	// 	os.Exit(1)
	// }
}
