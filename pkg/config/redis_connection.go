package config

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/redis/go-redis/v9"
)

var REDIS_DB_NUMBER = &env.Env.RedisDBNumber
var REDIS_PASSWORD = &env.Env.RedisPassword

func RedisConn() (*redis.Client, error) {
	redisCon, err := ConnectionURLBuilder("redis")
	if err != nil {
		return nil, err
	}
	options := &redis.Options{
		Addr:     redisCon,
		Password: *REDIS_PASSWORD,
		DB:       *REDIS_DB_NUMBER,
	}

	return redis.NewClient(options), nil
}
