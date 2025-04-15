package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topicCredentialCreated = "credential_creation"
	topicCredentialUpdated = "credential_update"
	topicCredentialDeleted = "credential_deletion"
	topicCredentialShared  = "credential_shared"
)

func main() {
	// Get Kafka broker address from environment variable or use default
	kafkaBroker := os.Getenv("KAFKA_BROKERS")
	if kafkaBroker == "" {
		kafkaBroker = "kafka:29092"
	}

	// Create a new Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Balancer: &kafka.LeastBytes{},
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that will be cancelled when we receive a signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a goroutine to handle signals
	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal. Cleaning up...")
		cancel()
	}()

	// Generate test messages
	fmt.Println("Starting to generate test messages...")

	// Create some test users and groups
	users := []string{"user-001", "user-002", "user-003"}
	groups := []string{"group-001", "group-002"}
	credentials := []string{"cred-001", "cred-002", "cred-003", "cred-004"}

	// Generate messages in a loop until context is cancelled
	for ctx.Err() == nil {
		// Randomly select a user, group, and credential
		userID := users[rand.Intn(len(users))]
		groupID := groups[rand.Intn(len(groups))]
		credentialID := credentials[rand.Intn(len(credentials))]

		// Randomly select a message type (now 4 types)
		messageType := rand.Intn(4)
		now := time.Now()

		var message interface{}
		var topic string
		var action string

		switch messageType {
		case 0:
			// Credential Created
			message = &CredentialCreated{
				UserID:       userID,
				GroupID:      groupID,
				CredentialID: credentialID,
				CreatedAt:    now,
			}
			topic = topicCredentialCreated
			action = "created"
		case 1:
			// Credential Updated
			message = &CredentialUpdated{
				UserID:       userID,
				GroupID:      groupID,
				CredentialID: credentialID,
				UpdatedAt:    now,
			}
			topic = topicCredentialUpdated
			action = "updated"
		case 2:
			// Credential Deleted
			message = &CredentialDeleted{
				UserID:       userID,
				GroupID:      groupID,
				CredentialID: credentialID,
				DeletedAt:    now,
			}
			topic = topicCredentialDeleted
			action = "deleted"
		case 3:
			// Credential Shared
			message = &CredentialShared{
				UserID:       userID,
				GroupID:      groupID,
				CredentialID: credentialID,
				SharedAt:     now,
				IPAddress:    fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
				UserAgent:    "Mozilla/5.0 (Test Browser)",
				IsOneTime:    rand.Intn(2) == 0,
			}
			topic = topicCredentialShared
			action = "shared"
		}

		// Marshal the message to JSON
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		// Create the Kafka message
		kafkaMessage := kafka.Message{
			Topic: topic,
			Key:   []byte(credentialID),
			Value: messageBytes,
		}

		// Write the message to Kafka
		err = writer.WriteMessages(ctx, kafkaMessage)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			continue
		}

		// Log the message
		log.Printf("Sent message: User %s %s credential %s in group %s",
			userID, action, credentialID, groupID)

		// Sleep for a random interval between 1 and 3 seconds
		time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
	}

	fmt.Println("Test producer stopped")
}
