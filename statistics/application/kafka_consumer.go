package application

import (
	"context"
	"log"

	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure/kafka"
	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure/repositories"
)

// KafkaConsumer represents the Kafka consumer application
type KafkaConsumer struct {
	consumer *kafka.Consumer
}

// NewKafkaConsumer creates a new Kafka consumer application
func NewKafkaConsumer(brokers []string, groupID string, repo *repositories.EventSourcingRepository) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(brokers, groupID, repo)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
	}, nil
}

// Setup initializes the Kafka consumer
func (k *KafkaConsumer) Setup() error {
	log.Println("Setting up Kafka consumer")
	return nil
}

// Ignite starts the Kafka consumer?
func (k *KafkaConsumer) Ignite() error {
	log.Println("Starting Kafka consumer")
	if err := k.consumer.Start(context.Background()); err != nil {
		return err
	}
	return nil
}

// Stop stops the Kafka consumer
func (k *KafkaConsumer) Stop() error {
	log.Println("Stopping Kafka consumer...")
	if err := k.consumer.Stop(context.Background()); err != nil {
		return err
	}
	return nil
}

// Shutdown implements the Application interface
func (k *KafkaConsumer) Shutdown() error {
	return k.Stop()
}
