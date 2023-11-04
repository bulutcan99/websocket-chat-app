package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ChatRoutes(r fiber.Router, cont *controller.HubController) {
	route := r.Group("/chat-ws")
	route.Post("/create-room", middleware.JWTProtection(), cont.CreateRoom)
	route.Get("/get-all-rooms", middleware.JWTProtection(), cont.GetAllRooms)
	route.Get("/get-available-rooms", middleware.JWTProtection(), cont.GetAvailableRooms)
	route.Get("/get-room-clients", middleware.JWTProtection(), cont.GetRoomClients)
	route.Use(middleware.FiberSocketUpgrade)
	route.Get("/join-room", middleware.JWTProtection(), websocket.New(cont.JoinRoom))
}
