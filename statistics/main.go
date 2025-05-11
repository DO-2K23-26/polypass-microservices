package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	esdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	"github.com/polypass/polypass-microservices/statistics/application/services"
	"github.com/polypass/polypass-microservices/statistics/config"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
	"github.com/polypass/polypass-microservices/statistics/infrastructure/api"
	"github.com/polypass/polypass-microservices/statistics/infrastructure/kafka"
	"github.com/polypass/polypass-microservices/statistics/infrastructure/repositories"
)

// No longer needed as we're using the config package

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context that listens for the interrupt signal from the OS
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up channel on which to send signal notifications
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine to listen for OS signals
	go func() {
		<-c
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Connect to EventStoreDB
	esdbClient, err := connectToEventStoreDB(ctx, cfg.EventStoreDB.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to EventStoreDB: %v", err)
	}
	defer esdbClient.Close()

	// Create repositories
	eventRepo := repositories.NewEventStoreDBEventRepository(esdbClient)
	metricRepo := repositories.NewEventStoreDBMetricRepository(esdbClient)

	// Create services
	eventService := services.NewEventService(eventRepo)

	// Create metric calculators
	calculators := []models.MetricCalculator{
		// Add your metric calculators here
		models.NewCredentialCountCalculator(),
		models.NewCredentialAccessCountCalculator(),
		// Add more calculators as needed
	}

	metricService := services.NewMetricService(metricRepo, eventService, calculators)

	// Create Kafka consumer
	kafkaConfig := kafka.EventConsumerConfig{
		BootstrapServers: cfg.Kafka.BootstrapServers,
		GroupID:          cfg.Kafka.GroupID,
		Topics:           cfg.Kafka.Topics,
	}

	eventConsumer, err := kafka.NewEventConsumer(kafkaConfig, eventService)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start Kafka consumer
	err = eventConsumer.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer eventConsumer.Stop()

	// Create HTTP router
	router := mux.NewRouter()

	// Create REST handler
	restHandler := api.NewRestHandler(metricService, eventService)
	restHandler.RegisterRoutes(router)

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.HTTP.ListenAddress,
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", cfg.HTTP.ListenAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for context cancellation (shutdown signal)
	<-ctx.Done()
	log.Println("Shutting down...")

	// Create a deadline for server shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}

// connectToEventStoreDB connects to EventStoreDB
func connectToEventStoreDB(ctx context.Context, connectionString string) (*esdb.Client, error) {
	// Create connection settings
	settings, err := esdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Create client
	client, err := esdb.NewClient(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create EventStoreDB client: %w", err)
	}

	log.Println("Connected to EventStoreDB")
	return client, nil
}
