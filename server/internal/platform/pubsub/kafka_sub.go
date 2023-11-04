package pubsub

import (
	"context"
	"github.com/IBM/sarama"
	config_kafka "github.com/bulutcan99/go-websocket/pkg/config/kafka"
	"go.uber.org/zap"
)

type ConsumerCallback func(data []byte)

type KafkaSubscriberInterface interface {
	Subscribe(topic string, consumerCallback ConsumerCallback) error
}

type KafkaSubscriber struct {
	Consumer sarama.Consumer
	Context  context.Context
}

func NewKafkaSubscriber(kafka config_kafka.ConsumerKafka) *KafkaSubscriber {
	return &KafkaSubscriber{
		Consumer: kafka.Consumer,
		Context:  kafka.Context,
	}
}

func (k *KafkaSubscriber) Subscribe(topic string, consumerCallback ConsumerCallback) error {
	partitionConsumer, err := k.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if nil != err {
		zap.S().Errorf("Error while consuming the Kafka topic: %s", err)
		return err
	}

	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		if nil != consumerCallback {
			consumerCallback(msg.Value)
		}
	}
}
