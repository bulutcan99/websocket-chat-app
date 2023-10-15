package config

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
)

var (
	DB_HOST     = &env.Env.DbHost
	DB_PORT     = &env.Env.DbPort
	DB_USER     = &env.Env.DbUser
	DB_PASSWORD = &env.Env.DbPassword
	DB_NAME     = &env.Env.DbName
	DB_SSL_MODE = &env.Env.DbSSLMode
	REDIS_HOST  = &env.Env.RedisHost
	REDIS_PORT  = &env.Env.RedisPort
	SERVER_HOST = &env.Env.ServerHost
	SERVER_PORT = &env.Env.ServerPort
)

func ConnectionURLBuilder(n string) (string, error) {
	var url string
	switch n {
	case "postgres":
		url = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			*DB_HOST,
			*DB_PORT,
			*DB_USER,
			*DB_PASSWORD,
			*DB_NAME,
			*DB_SSL_MODE,
		)
	case "redis":
		url = fmt.Sprintf(
			"%s:%d",
			*REDIS_HOST,
			*REDIS_PORT,
		)
	case "fiber":
		url = fmt.Sprintf(
			"%s:%d",
			*SERVER_HOST,
			*SERVER_PORT,
		)
	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}
