package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes Base Route = /v1/auth
func AuthRoutes(r fiber.Router, au *controller.AuthController) {
	route := r.Group("/auth")
	route.Post("/register", au.UserRegister)
	route.Post("/login", au.UserLogin)
	route.Post("/logout", middleware.JWTProtected(), au.UserLogOut)
}
