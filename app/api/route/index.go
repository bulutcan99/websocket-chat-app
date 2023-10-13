package route

import (
	"github.com/bulutcan99/go-websocket/pkg/index"
	"github.com/gofiber/fiber/v2"
)

// IndexRoutes Base Route = /MSN/
func IndexRoutes(b string, app *fiber.App) {
	app.Get(b, index.Homepage)
	app.Get(b+"signup", index.Signup)
	app.Get(b+"login", index.Login)
	app.Get(b+"chats", index.Chats)
}
