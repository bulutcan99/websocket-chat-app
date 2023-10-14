package index

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/gofiber/fiber/v2"
)

var (
	DB_HOST = &env.Env.DbHost
	DB_PORT = &env.Env.DbPort
)

func Homepage(c *fiber.Ctx) error {
	return c.Render("index",
		fiber.Map{
			"Awesome": "Go-Chat",
			"Home":    "active",
		})
}

func Signup(c *fiber.Ctx) error {
	return c.Render("sign-up-screen",
		fiber.Map{
			"Awesome":  "GO-WS-CHAT",
			"Register": "active",
		})
}

func Login(c *fiber.Ctx) error {
	return c.Render("login-screen",
		fiber.Map{
			"Awesome": "Go-Chat",
			"Login":   "active",
		})
}

func Chats(c *fiber.Ctx) error {
	return c.Render("chats-screen",
		fiber.Map{
			"HOST":  *DB_HOST,
			"PORT":  *DB_PORT,
			"Chats": "active",
		})
}
