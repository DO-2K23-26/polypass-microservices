package producer

// The only purpose of this file is debug. 
// If you need to populate a topic with some messages feel free to use this file.


import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/authz-service/infrastructure"
)

type Producer struct {
	kafka infrastructure.KafkaAdapter
}

func NewProducer(kafka infrastructure.KafkaAdapter) *Producer {
	return &Producer{
		kafka: kafka,
	}
}

func (p *Producer) StartPeriodicProducer(topic string, ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping periodic producer...")
			return
		case t := <-ticker.C:
			message := fmt.Sprintf("Message produced at %s", t.Format(time.RFC3339))
			err := p.kafka.Produce(topic, []byte(message))
			if err != nil {
				fmt.Printf("Error producing message: %v\n", err)
			} else {
				log.Printf("Message produced: %s\n", message)
			}
		}
	}
}
