package middleware

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/pkg/env"
	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var (
	JWT_SECRET_KEY = &env.Env.JwtSecretKey
)

func JWTProtection() func(*fiber.Ctx) error {
	config := jwtMiddleware.Config{
		SigningKey: jwtMiddleware.SigningKey{
			JWTAlg: jwtMiddleware.HS256,
			Key:    []byte(*JWT_SECRET_KEY)},
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

	fmt.Println("ERROR: ", err)
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
