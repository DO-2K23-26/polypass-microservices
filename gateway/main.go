package main

import (
	"os"

	"github.com/DO-2K23-26/polypass-microservices/gateway/application/graphql"
	"github.com/DO-2K23-26/polypass-microservices/gateway/application/searchhttp"
	"github.com/DO-2K23-26/polypass-microservices/gateway/config"
	"github.com/DO-2K23-26/polypass-microservices/gateway/core"
	"github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure/organizations"
	"github.com/DO-2K23-26/polypass-microservices/gateway/infrastructure/search"
	"github.com/optique-dev/optique"
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
	cycle := NewCycle()

	// Initialize organization REST client
	organizations_api := organizations.New(conf.OrganizationConfig)
	cycle.AddRepository(organizations_api)

	// Initialize gRPC search client
	search_api := search.New(conf.SearchConfig)
	// immediately dial search service so client is non-nil
	if err := search_api.Setup(); err != nil {
		optique.Error("cannot connect to search service: " + err.Error())
		os.Exit(1)
	}
	cycle.AddRepository(search_api)

	// add services
	organizations_service := core.NewOrganizationsService(organizations_api)
	// now that search_api is Setup, create core search service
	search_service := core.NewSearchService(search_api.SearchServiceClient())

	// add applications
	api_gateway := graphql.NewHttp(conf.GraphQL)
	api_gateway.WithHandler(graphql.NewHealthController())
	// REST search endpoints
	searchController := searchhttp.NewController(search_service)
	api_gateway.WithHandler(searchController)
	api_gateway.WithHandler(graphql.NewGraphQL(organizations_service))

	cycle.AddApplication(api_gateway)

	if conf.Bootstrap {
		err := cycle.Setup()
		if err != nil {
			optique.Error(err.Error())
			cycle.Stop()
			os.Exit(1)
		}
	}

	err = cycle.Ignite()
}
