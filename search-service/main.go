package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DO-2K23-26/polypass-microservices/search-service/app"
	// "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/credential"
	"github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"

	// "github.com/DO-2K23-26/polypass-microservices/search-service/repositories/tags"

	// "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
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
	tagRepo := *app.TagRepository
	credRepo := *app.CredentialRepository

	// credRepo.Create(credential.CreateCredentialQuery{
	// 	ID:    "cred1",
	// 	Title: "test cred",
	// 	Tags: []types.Tag{{
	// 		ID:       "tagtest1",
	// 		Name:     "Deletable tag",
	// 		FolderId: "folder",
	// 	}},
	// })

	// _, err = tagRepo.Create(tags.CreateTagQuery{
	// 	ID:       "tagtest1",
	// 	Name:     "Deletable tag",
	// 	FolderID: "folder",
	// })

	err = tagRepo.Delete(tags.DeleteTagQuery{ID: "tagtest1"})
	if err != nil {
		panic(err)
	}

	res, err := credRepo.Get(credential.GetCredentialQuery{ID: "cred1"})
	if err != nil {
		panic(err)
	} else {
		log.Println(res)
	}

	restag, err := tagRepo.Get(tags.GetTagQuery{ID: "tagtest1"})
	if err != nil {
		panic(err)
	} else {
		log.Println(restag.Tag)
	}
	// else {
	// 	log.Println(res)
	// }

	// res,err =test.Get( tags.GetTagQuery{
	// 	ID: "tagtest",
	// })
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	log.Println(res)
	// }
	// err = app.Start()
	// if err != nil {
	// 	fmt.Printf("Failed to start the application: %v\n", err)
	// 	os.Exit(1)
	// }
}
