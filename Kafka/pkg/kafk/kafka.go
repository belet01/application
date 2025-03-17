package kafk

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer(producerConfig kafka.ConfigMap) *kafka.Producer {
	producer, err := kafka.NewProducer(&producerConfig)
	if err != nil {
		return producer
	}
	return producer
}

func NewConsumer(consumerConfig kafka.ConfigMap) (*kafka.Consumer, error) {

	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {

		return nil, fmt.Errorf("consumer oluşturulamadı: %v", err)
	}
	return consumer, nil
}
