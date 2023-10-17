package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/gofiber/fiber/v2"
)

// AuthRoutes Base Route = /v1/auth
func AuthRoutes(r fiber.Router, authController *controller.AuthController) {
	route := r.Group("/v1/auth")
	route.Post("/register", authController.UserSignUp)
	route.Post("/login", authController.UserSignIn)
}
