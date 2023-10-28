package config_redis

import (
	"context"
	config_builder "github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
)

var (
	doOnce          sync.Once
	client          *redis.Client
	REDIS_DB_NUMBER = &env.Env.RedisDBNumber
	REDIS_PASSWORD  = &env.Env.RedisPassword
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

func NewRedisConnection() *Redis {
	ctx := context.Background()

	redisCon, err := config_builder.ConnectionURLBuilder("redis")
	if err != nil {
		panic(err)
	}

	options := &redis.Options{
		Addr:     redisCon,
		Password: *REDIS_PASSWORD,
		DB:       *REDIS_DB_NUMBER,
	}

	doOnce.Do(func() {
		redisClient := redis.NewClient(options)
		_, err = redisClient.Ping(ctx).Result()
		if err != nil {
			panic(err)
		}

		redisClient.Ping(ctx)

		client = redisClient
	})

	zap.S().Info("Connected to Redis successfully.")
	return &Redis{
		Client:  client,
		Context: ctx,
	}
}

func (r *Redis) Close() {
	if err := r.Client.Close(); err != nil {
		zap.S().Errorf("Error while closing the Redis connection: %s", err)
	}

	zap.S().Info("Connection to Redis closed successfully")
}
