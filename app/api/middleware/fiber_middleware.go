package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func getCorsConfig() *cors.Config {
	return &cors.Config{
		AllowCredentials: true,
	}
}

func MiddlewareFiber(m *fiber.App) {
	m.Use(
		cors.New(*getCorsConfig()),
		logger.New(),
	)
}

func GetNextMiddleWare(c *fiber.Ctx) error {
	return c.Next()
}
