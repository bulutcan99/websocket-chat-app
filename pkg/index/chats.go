package index

import (
	"github.com/bulutcan99/go-websocket/pkg/env"

	"github.com/gofiber/fiber/v2"
)

func Chats(c *fiber.Ctx) error {
	EnvStatus := false
	Env := env.ParseEnv()
	if Env.StageStatus == "dev" {
		EnvStatus = true
	}

	return c.Render("chats-screen",
		fiber.Map{
			"HOST":  Env.DbHost,
			"PORT":  Env.DbPort,
			"Chats": "active",
			"Env":   EnvStatus,
		})
}
