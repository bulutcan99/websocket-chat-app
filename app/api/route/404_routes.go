package route

import (
	"github.com/bulutcan99/go-websocket/model"
	"github.com/gofiber/fiber/v2"
)

func NotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(
				model.Error{
					Message: "404 Not Found",
					Error:   true,
				})
		},
	)
}
