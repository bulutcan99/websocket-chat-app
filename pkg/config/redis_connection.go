package config

import (
	"context"
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/redis/go-redis/v9"
)

var REDIS_DB_NUMBER = &env.Env.RedisDBNumber
var REDIS_PASSWORD = &env.Env.RedisPassword

type Redis struct {
	Client  *redis.Client
	Context context.Context // Bağlam ekleniyor
}

func NewRedisConnection() *Redis {
	redisCon, err := ConnectionURLBuilder("redis")
	if err != nil {
		panic(err)
	}

	options := &redis.Options{
		Addr:     redisCon,
		Password: *REDIS_PASSWORD,
		DB:       *REDIS_DB_NUMBER,
	}

	client := redis.NewClient(options)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &Redis{Client: client, Context: context.Background()}
}

func (r *Redis) Close() {
	if err := r.Client.Close(); err != nil {
		fmt.Printf("Error while closing the Redis connection: %s\n", err)
	}
}
