package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	domain "github.com/DO-2K23-26/polypass-microservices/statistics/domain/repositories"
	"github.com/DO-2K23-26/polypass-microservices/statistics/infrastructure/repositories"
	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	reader     *kafka.Reader
	repo       *repositories.EventSourcingRepository
	topics     []string
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, repo *repositories.EventSourcingRepository) (*Consumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}
	if groupID == "" {
		return nil, fmt.Errorf("group ID is required")
	}
	if repo == nil {
		return nil, fmt.Errorf("repository is required")
	}

	topics := []string{
		"credential_creation",
		"credential_update",
		"credential_deletion",
		"credential_shared",
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
		StartOffset: kafka.FirstOffset,
		MaxWait:     10 * time.Second,
		Logger:      kafka.LoggerFunc(log.Printf),
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		reader:     reader,
		repo:       repo,
		topics:     topics,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.consumeMessages()
	}()

	return nil
}

// Stop gracefully shuts down the consumer
func (c *Consumer) Stop(ctx context.Context) error {
	c.cancelFunc()
	c.wg.Wait()
	return c.reader.Close()
}

// consumeMessages consumes messages from all configured topics
func (c *Consumer) consumeMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(c.ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var eventType domain.EventType
			switch msg.Topic {
			case "credential_creation", "credential_update", "credential_deletion":
				eventType = repositories.CredentialEvent
			case "credential_shared":
				eventType = repositories.AccessEvent
			default:
				log.Printf("Unknown topic: %s", msg.Topic)
				continue
			}

			if err := c.repo.StoreEvent(c.ctx, domain.Event{
				ID:        string(msg.Key),
				Type:      eventType,
				Data:      msg.Value,
				Timestamp: time.Now(),
			}); err != nil {
				log.Printf("Error storing event: %v", err)
				continue
			}
		}
	}
}
