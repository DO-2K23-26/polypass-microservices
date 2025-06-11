package infrastructure

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func (k *KafkaAdapter) Consume(topic string, handleMessage func(*kafka.Message) error, handleError func(error)) error {
	println("Consuming topic:", topic)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.host,
		"group.id":          k.clientId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		println("Error creating consumer:", err)
		return err
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	
	run := true
	for run == true {
		select {
        case sig := <-sigchan:
            fmt.Printf("Caught signal %v: terminating\n", sig)
            run = false
        default:
            ev, err := consumer.ReadMessage(-1)
            if err != nil {
                if handleError != nil {
					handleError(err)
				}
                continue
            }
            fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
                *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			if handleMessage != nil {
				err := handleMessage(ev)
				if err != nil {
					if handleError != nil {
						handleError(err)
					}
					continue
				}
			}
			
			_, err = k.consumer.CommitMessage(ev)
			if err != nil && handleError != nil {
				handleError(err)
			}
        }

		if err != nil {
			
			continue
		}
	}
	return err
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
