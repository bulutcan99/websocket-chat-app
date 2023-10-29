package token

import (
	"github.com/bulutcan99/go-websocket/pkg/env"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_KEY = &env.Env.JwtSecretKey
)

func ExtractTokenMetaData(c *fiber.Ctx) (*TokenMetaData, error) {
	token, err := verifyToken(c)
	if err != nil {
		_ = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		_ = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil, err
	}

	if err != nil {
		_ = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
		return nil, err
	}

	expires := int64(claims["exp"].(float64))
	role := claims["role"].(string)
	id := claims["id"].(string)
	email := claims["email"].(string)
	return &TokenMetaData{
		UUID:    id,
		Expires: expires,
		Role:    role,
		Email:   email,
	}, nil
}

func ExtractToken(c *fiber.Ctx) string {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	jwtKey := []byte(*JWT_KEY)
	return jwtKey, nil
}
