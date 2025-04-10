package main

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/config"
	"github.com/elastic/go-elasticsearch/v8"
)

// @title Optique application TO CHANGE
// @version 1.0
// @description This is a sample application
// @contact.name Courtcircuits
// @contact.url https://github.com/Courtcircuits
// @contact.email tristan-mihai.radulescu@etu.umontpellier.fr
func main() {
	esConfig := elasticsearch.Config{Addresses: []string{"http://localhost:9200"}, Password: "search-password"}
	es, _ := elasticsearch.NewClient(esConfig)
	esPingRes, err := es.Ping()
	if err != nil {
		println(err.Error())
		return
	}

	if esPingRes != nil {

		println(esPingRes.String())
	}
	conf, err := config.LoadConfig()
	if err != nil {
		config.HandleError(err)
	}
	println(conf)
}
