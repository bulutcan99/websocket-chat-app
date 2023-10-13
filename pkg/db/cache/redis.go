package cache

import (
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/redis/go-redis/v9"
)

func RedisConn() (*redis.Client, error) {
	Env := env.ParseEnv()
	redisDbNumber := Env.RedisDBNumber
	redisCon, err := config.ConnectionURLBuilder("redis")
	if err != nil {
		return nil, err
	}
	redisPass := Env.RedisPassword
	options := &redis.Options{
		Addr:     redisCon,
		Password: redisPass,
		DB:       redisDbNumber,
	}

	return redis.NewClient(options), nil
}
