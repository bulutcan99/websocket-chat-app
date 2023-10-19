package middleware

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	JWT_SECRET_KEY = &env.Env.JwtSecretKey
)

func JWTProtected() func(*fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(*JWT_SECRET_KEY)},
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
