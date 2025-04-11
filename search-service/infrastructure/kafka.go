package infrastructure

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaAdapter struct {
	host     string
	clientId string
	producer *kafka.Producer
	consumer *kafka.Consumer
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
	return &KafkaAdapter{
		host:     host,
		clientId: clientId,
		producer: producer,
		consumer: consumer,
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

func (k *KafkaAdapter) Consume(topic string, handleMessage func(*kafka.Message) error, handleError func(error)) error {
	err := k.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			if handleError != nil {
				handleError(err)
			}
			continue
		}
		if handleMessage != nil {
			err := handleMessage(msg)
			if err != nil {
				if handleError != nil {
					handleError(err)
				}
				continue
			}
		}

		_, err = k.consumer.CommitMessage(msg)
		if err != nil && handleError != nil {
			handleError(err)
		}
	}
}
