package events

import "github.com/confluentinc/confluent-kafka-go/kafka"

type IUserEventController interface {
    AddUserToFolder(*kafka.Message) error
    RemoveUserToFolder(*kafka.Message) error
}
