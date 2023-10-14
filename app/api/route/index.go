package route

import (
	"github.com/bulutcan99/go-websocket/app/index"
	"github.com/gofiber/fiber/v2"
)

func Index(b string, app *fiber.App) {
	app.Get(b, index.Homepage)
	app.Get(b+"register", index.Signup)
	app.Get(b+"login", index.Login)
	app.Get(b+"chats", index.Chats)
}
