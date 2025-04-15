package main

import (
	"log"
	"os"
	"strings"

	"github.com/DO-2K23-26/polypass-microservices/statistics/application"
	"github.com/DO-2K23-26/polypass-microservices/statistics/application/services"
	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure/repositories"
)

// @title Statistics service
// @version 1.0
// @description This application is used to collect statistics about the usage of the credentials in the application
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Get configuration from environment variables
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "statistics-group")
	httpPort := getEnv("HTTP_PORT", "8080")

	// Create cycle
	c := NewCycle()

	// Initialize repository
	repo := repositories.NewEventSourcingRepository()
	c.AddRepository(repo)

	// Initialize service
	service := services.NewMetricsService(repo)

	// Create HTTP API application
	httpAPI := application.NewHTTPAPI(httpPort, service)
	c.AddApplication(httpAPI)

	// Create Kafka consumer application
	kafkaConsumer, err := application.NewKafkaConsumer(kafkaBrokers, kafkaGroupID, repo)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	c.AddApplication(kafkaConsumer)

	// Run cycle
	if err := c.Setup(); err != nil {
		log.Fatalf("Failed to setup application: %v", err)
	}

	if err := c.Ignite(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
