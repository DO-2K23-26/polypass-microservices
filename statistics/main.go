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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	esdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	"github.com/polypass/polypass-microservices/statistics/application/services"
	"github.com/polypass/polypass-microservices/statistics/config"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
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

	// Create services
	eventService := services.NewEventService(eventRepo)

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
	r := gin.Default()

	// Configuration CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/statistics/analyze", func(c *gin.Context) {
		var credentials []map[string]interface{}
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Analyse des mots de passe
		weakPasswords := make(map[string][]string)
		strongPasswords := make(map[string][]string)
		reusedPasswords := make(map[string][]string)
		oldPasswords := make(map[string][]string)
		breachedPasswords := make(map[string][]string)

		// Map pour suivre les mots de passe réutilisés
		passwordOccurrences := make(map[string][]string)

		// Date limite pour les mots de passe anciens (1 an)
		oneYearAgo := time.Now().AddDate(-1, 0, 0)

		for _, cred := range credentials {
			password, ok := cred["password"].(string)
			if !ok {
				continue
			}

			id, ok := cred["id"].(string)
			if !ok {
				continue
			}

			// Vérifier si le mot de passe est fort
			isStrong := models.IsStrongPassword(password)

			if isStrong {
				strongPasswords[password] = append(strongPasswords[password], id)
			} else {
				weakPasswords[password] = append(weakPasswords[password], id)
			}

			// Suivre les occurrences de chaque mot de passe
			passwordOccurrences[password] = append(passwordOccurrences[password], id)

			// Vérifier si le mot de passe est ancien
			if lastUpdated, ok := cred["lastUpdated"].(string); ok {
				lastUpdatedTime, err := time.Parse("2006-01-02", lastUpdated)
				if err == nil && lastUpdatedTime.Before(oneYearAgo) {
					oldPasswords[password] = append(oldPasswords[password], id)
				}
			}

			// Vérifier si le mot de passe est compromis
			if isBreached, err := models.CheckPasswordBreach(password); err == nil && isBreached {
				breachedPasswords[password] = append(breachedPasswords[password], id)
			}
		}

		// Identifier les mots de passe réutilisés
		for password, ids := range passwordOccurrences {
			if len(ids) > 1 {
				reusedPasswords[password] = ids
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"weakPasswords":     weakPasswords,
			"strongPasswords":   strongPasswords,
			"reusedPasswords":   reusedPasswords,
			"oldPasswords":      oldPasswords,
			"breachedPasswords": breachedPasswords,
		})
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    cfg.HTTP.ListenAddress,
		Handler: r,
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
