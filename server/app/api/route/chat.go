package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ChatRoutes(r fiber.Router, cont controller.HubController) {
	route := r.Group("/chat-ws")
	route.Post("/create-room", middleware.JWTProtection(), cont.CreateRoom)
	route.Use(middleware.FiberSocketUpgrade)
	route.Post("/join-room", middleware.JWTProtection(), websocket.New(cont.JoinRoom))
}
