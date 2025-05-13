package infrastructure

import (
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/avro-schemas/schemautils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaAdapter struct {
	host     string
	clientId string
	producer *kafka.Producer
	consumer *kafka.Consumer
	admin    *kafka.AdminClient
}

func NewKafkaAdapter(host string, clientId string) (*KafkaAdapter, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": host,
		"client.id":         clientId,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": host,
		"group.id":          clientId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	admin, err := kafka.NewAdminClient(config)
	if err != nil {
		return nil, err
	}
	return &KafkaAdapter{
		host:     host,
		clientId: clientId,
		producer: producer,
		consumer: consumer,
		admin:    admin,
	}, nil
}

func (k *KafkaAdapter) Produce(topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	msg := e.(*kafka.Message)

	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func (k *KafkaAdapter) CheckHealth() bool {
	// Use the AdminClient to check the status of the brokers
	_, err := k.admin.GetMetadata(nil, true, 5000)
	if err != nil {
		fmt.Println("Error checking Kafka health:", err)
		return false
	}
	return true
}

func (k *KafkaAdapter) ProduceAvro(topic string, encoder *schemautils.AvroEncoder, data map[string]interface{}) error {
	value, err := encoder.Encode(data)
	if err != nil {
		return err
	}
	return k.Produce(topic, value)
}
