package config_kafka

import (
	"context"
	config_builder "github.com/bulutcan99/go-websocket/pkg/config"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var (
	kafkaProducer sarama.AsyncProducer
)

type ProducerKafka struct {
	Producer sarama.AsyncProducer
	Context  context.Context
}

func NewKafkaProducerConnection() *ProducerKafka {
	ctx := context.Background()
	kafkaCon, err := config_builder.ConnectionURLBuilder("kafka")
	if err != nil {
		panic(err)
	}

	kafkaConSplit := strings.Split(kafkaCon, ":")
	config := sarama.NewConfig()
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	doOnce.Do(func() {
		kafkaClient, err = sarama.NewClient(kafkaConSplit, config)
		if err != nil {
			panic(err)
		}

		if err := kafkaClient.RefreshMetadata(); err != nil {
			panic(err)
		}

		producer, err := sarama.NewAsyncProducerFromClient(kafkaClient)
		if nil != err {
			zap.S().Errorf("Error while creating the Kafka Producer connection: %s", err)
		}

		kafkaProducer = producer
	})

	zap.S().Info("Connected to Kafka Producer successfully.")
	return &ProducerKafka{
		Producer: kafkaProducer,
		Context:  ctx,
	}
}

func (k *ProducerKafka) Close() {
	if err := k.Producer.Close(); err != nil {
		zap.S().Errorf("Error while closing the Kafka Producer connection: %s\n", err)
	}

	zap.S().Info("Connection to Kafka Producer closed successfully")
}
