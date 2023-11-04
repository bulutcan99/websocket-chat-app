package pubsub

import (
	"context"
	"github.com/IBM/sarama"
	config_kafka "github.com/bulutcan99/go-websocket/pkg/config/kafka"
)

type KafkaPublisherInterface interface {
	Publish(topic string, message []byte) error
}

type KafkaPublisher struct {
	Producer sarama.AsyncProducer
	Context  context.Context
}

func NewKafkaPublisher(kafka config_kafka.ProducerKafka) *KafkaPublisher {
	return &KafkaPublisher{
		Producer: kafka.Producer,
		Context:  kafka.Context,
	}
}

func (k *KafkaPublisher) Publish(topic string, message []byte) error {
	be := sarama.ByteEncoder(message)
	k.Producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: be}
	return nil
}
