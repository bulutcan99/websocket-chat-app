package config

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"time"

	"github.com/gofiber/fiber/v2"
)

var READ_TIMEOUT_SECONDS_COUNT = &env.Env.ServerReadTimeout

func FiberConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(*READ_TIMEOUT_SECONDS_COUNT),
	}
}
