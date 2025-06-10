package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/polypass/polypass-microservices/statistics/application/services"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
)

// EventConsumer consumes events from Kafka topics
type EventConsumer struct {
	consumer     sarama.ConsumerGroup
	eventService *services.EventService
	topics       []string
	running      bool
	mutex        sync.Mutex
	ready        chan bool
}

// EventConsumerConfig holds configuration for the Kafka consumer
type EventConsumerConfig struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
}

// consumerGroupHandler implements the sarama.ConsumerGroupHandler interface
type consumerGroupHandler struct {
	eventService *services.EventService
	ready        chan bool
}

// NewEventConsumer creates a new instance of EventConsumer
func NewEventConsumer(config EventConsumerConfig, eventService *services.EventService) (*EventConsumer, error) {
	// Create Sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Return.Errors = true

	// Split bootstrap servers string into slice
	bootstrapServers := strings.Split(config.BootstrapServers, ",")

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(
		bootstrapServers,
		config.GroupID,
		saramaConfig,
	)
	if err != nil {
		return nil, err
	}

	return &EventConsumer{
		consumer:     consumer,
		eventService: eventService,
		topics:       config.Topics,
		running:      false,
		ready:        make(chan bool),
	}, nil
}

// Start begins consuming events from Kafka
func (c *EventConsumer) Start(ctx context.Context) error {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return nil
	}
	c.running = true
	c.mutex.Unlock()

	// Start consuming in a separate goroutine
	go func() {
		handler := &consumerGroupHandler{
			eventService: c.eventService,
			ready:        c.ready,
		}

		// Consume until context is cancelled
		for c.isRunning() {
			// Consume from topics
			if err := c.consumer.Consume(ctx, c.topics, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}

			// Check if context was cancelled, signaling a shutdown
			if ctx.Err() != nil {
				return
			}

			// Reset ready channel for next consume
			c.ready = make(chan bool)
		}
	}()

	// Wait until the consumer is ready
	<-c.ready
	log.Println("Kafka consumer started")

	return nil
}

// Stop stops consuming events from Kafka
func (c *EventConsumer) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.running {
		return nil
	}

	c.running = false
	return c.consumer.Close()
}

// isRunning returns whether the consumer is running
func (c *EventConsumer) isRunning() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.running
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Process messages
	for message := range claim.Messages() {
		if err := h.processMessage(session.Context(), message); err != nil {
			log.Printf("Error processing message: %v", err)
		}
		// Mark message as processed
		session.MarkMessage(message, "")
	}
	return nil
}

// processMessage processes a single Kafka message
func (h *consumerGroupHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event models.Event

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}

	// Set the timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Store the event
	return h.eventService.StoreEvent(ctx, event)
}
