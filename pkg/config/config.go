package config

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	Env := env.ParseEnv()
	readTimeoutSecondsCount := Env.ServerReadTimeout

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
