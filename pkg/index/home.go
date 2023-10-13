package index

import (
	"github.com/gofiber/fiber/v2"
)

func Homepage(c *fiber.Ctx) error {
	return c.Render("index-screen",
		fiber.Map{
			"Awesome": "GO-WS-CHAT",
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
