package config

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/gofiber/template/html/v2"
	"time"

	"github.com/gofiber/fiber/v2"
)

var READ_TIMEOUT_SECONDS_COUNT = &env.Env.ServerReadTimeout

func ConfigFiber() fiber.Config {
	return fiber.Config{
		ReadTimeout:  time.Second * time.Duration(*READ_TIMEOUT_SECONDS_COUNT),
		Prefork:      false, // Disable Prefork means that need only 1 instance running
		ServerHeader: "iMon",
		AppName:      "Go-Chat",
		Immutable:    true,
		Views:        getHandler(),
	}
}

func getHandler() *html.Engine {
	handler := html.New("./views", ".html")
	return handler
}
