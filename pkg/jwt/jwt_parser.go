package jwt

import (
	"github.com/bulutcan99/go-websocket/app/model"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_KEY = &env.Env.StageStatus
)

var jwtKey = []byte(*JWT_KEY)

func ExtractTokenMetaData(c *fiber.Ctx) (*model.TokenMetaData, error) {
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

	expires := int64(claims["expires"].(float64))
	role := claims["role"].(string)
	return &model.TokenMetaData{
		Expires: expires,
		Role:    role,
	}, nil
}

func extractToken(c *fiber.Ctx) string {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(jwtKey), nil
}
