package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes Base Route = /v1/auth
func AuthRoutes(r fiber.Router) {
	route := r.Group("/v1", middleware.GetNextMiddleWare)
	route.Post("/register", controller.UserSignUp)
}
