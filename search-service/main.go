package main

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
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
	print(conf)
}
