package route

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/gofiber/fiber/v2"
)

// UserRoutes Base Route = /user
func UserRoutes(r fiber.Router, user *controller.UserController) {
	// route := r.Group("/user")
	// route.Get("/self", middleware.JWTProtection(), user.GetUserSelf())
}
