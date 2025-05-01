package events

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ICredentialEventController interface {
	Create(*kafka.Message) error
	Update(*kafka.Message) error
	Delete(*kafka.Message) error
}
