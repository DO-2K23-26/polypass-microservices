package application

import (
	"context"
	"log"
	"net/http"

	"github.com/DO-2K23-26/polypass-microservices/statistics/application/services"
	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure/api"
	"github.com/gorilla/mux"
)

// HTTPAPI represents the HTTP API application
type HTTPAPI struct {
	server  *http.Server
	service *services.MetricsService
}

// NewHTTPAPI creates a new HTTP API application
func NewHTTPAPI(port string, service *services.MetricsService) *HTTPAPI {
	handler := api.NewHandler(service)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	return &HTTPAPI{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		service: service,
	}
}

// Setup initializes the HTTP API
func (a *HTTPAPI) Setup() error {
	log.Println("Setting up HTTP API")
	return nil
}

// Ignite starts the HTTP API
func (a *HTTPAPI) Ignite() error {
	log.Printf("Starting HTTP server on %s", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the HTTP API
func (a *HTTPAPI) Stop() error {
	log.Println("Stopping HTTP API...")
	if err := a.server.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

// Shutdown implements the Application interface
func (a *HTTPAPI) Shutdown() error {
	return a.Stop()
}
