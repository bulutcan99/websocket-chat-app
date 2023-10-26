package config_rabbitMq

import (
	"context"
	config_builder "github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"sync"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	doOnce       sync.Once
	rabbitMQ     RabbitMQ
	RABBITMQ_URL = &env.Env.RabbitMQUrl
)

type RabbitMQ struct {
	Client  *amqp.Connection
	Context context.Context
}

func NewRabbitMq() *RabbitMQ {
	ctx := context.Background()
	rabbitCon, err := config_builder.ConnectionURLBuilder("rabbitmq")
	if err != nil {
		panic(err)
	}

	doOnce.Do(func() {
		conn := RabbitMQConnection(rabbitCon)
		rabbitMQ = RabbitMQ{Client: conn, Context: ctx}
	})

	return &rabbitMQ
}

func RabbitMQConnection(rabbitMQUrl string) *amqp.Connection {
	conn, err := amqp.Dial(rabbitMQUrl)

	if err != nil {
		zap.S().Error("Failed to connect to RabbitMQ: ", err.Error())
	}

	zap.S().Infof("Connected to RabbitMQ successfully URL: %s", rabbitMQUrl)

	return conn
}

func (r *RabbitMQ) RabbitMQChannel() *amqp.Channel {
	ch, err := r.Client.Channel()

	if err != nil {
		zap.S().Error("Failed to open a channel: ", err.Error())
	}

	return ch
}

func (r *RabbitMQ) Close() {
	isClosed := r.Client.IsClosed()

	if !isClosed {
		err := r.Client.Close()
		if err != nil {
			zap.S().Errorf("Error while disconnecting from RabbitMQ: %s", err)
		}
	}

	zap.S().Infof("Connection to RabbitMQ closed successfully")
}
