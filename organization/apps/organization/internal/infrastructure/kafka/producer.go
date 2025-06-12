package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerAdapter struct {
	producer *kafka.Producer
}

func NewProducerAdapter(bootstrapServers, clientId string) (*ProducerAdapter, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         clientId,
	})
	if err != nil {
		return nil, err
	}

	return &ProducerAdapter{producer: p}, nil
}

func (p *ProducerAdapter) Publish(topic string, value []byte) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}
