package config_kafka

import (
	"context"
	"github.com/IBM/sarama"
	config_builder "github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var (
	doOnce        sync.Once
	kafkaClient   sarama.Client
	kafkaConsumer sarama.Consumer
	topic         = &env.Env.KafkaMessageTopic
)

type ConsumerKafka struct {
	Consumer sarama.Consumer
	Context  context.Context
}

func NewKafkaConsumerConnection() *ConsumerKafka {
	ctx := context.Background()
	kafkaCon, err := config_builder.ConnectionURLBuilder("kafka")
	if err != nil {
		panic(err)
	}

	kafkaConSplit := strings.Split(kafkaCon, ":")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	doOnce.Do(func() {
		kafkaClient, err = sarama.NewClient(kafkaConSplit, config)
		if err != nil {
			panic(err)
		}

		if err := kafkaClient.RefreshMetadata(); err != nil {
			panic(err)
		}

		consumer, err := sarama.NewConsumerFromClient(kafkaClient)
		if nil != err {
			zap.S().Errorf("Error while creating the Kafka Consumer connection: %s", err)
		}

		kafkaConsumer = consumer
	})

	zap.S().Info("Connected to Kafka Consumer successfully.")
	return &ConsumerKafka{
		Consumer: kafkaConsumer,
		Context:  ctx,
	}
}

func (k *ConsumerKafka) Close() {
	if err := k.Consumer.Close(); err != nil {
		zap.S().Errorf("Error while closing the Kafka Consumer connection: %s\n", err)
	}

	zap.S().Info("Connection to Kafka Consumer closed successfully")
}
