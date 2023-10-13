package config

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
)

func ConnectionURLBuilder(n string) (string, error) {
	Env := env.ParseEnv()
	var url string

	switch n {
	case "postgres":
		url = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			Env.DbHost,
			Env.DbPort,
			Env.DbUser,
			Env.DbPassword,
			Env.DbName,
			Env.DbSSLMode,
		)
	case "redis":
		url = fmt.Sprintf(
			"%s:%d",
			Env.RedisHost,
			Env.RedisPort,
		)
	case "fiber":
		url = fmt.Sprintf(
			"%s:%d",
			Env.ServerHost,
			Env.ServerPort,
		)
	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}
