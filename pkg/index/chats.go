package index

import (
	"github.com/bulutcan99/go-websocket/pkg/env"

	"github.com/gofiber/fiber/v2"
)

func Chats(c *fiber.Ctx) error {
	return c.Render("chats-screen",
		fiber.Map{
			"HOST":  &env.Env.DbHost,
			"PORT":  &env.Env.DbPort,
			"Chats": "active",
		})
}
