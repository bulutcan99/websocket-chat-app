package jwt

import (
	"github.com/google/uuid"
	"os"
	"strings"
	"websocket_chat/app/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func ExtractTokenMetaData(c *fiber.Ctx) (*model.TokenMetaData, error) {
	token, err := verifyToken(extractToken(c))
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	userID, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	expires := int64(claims["expires"].(float64))
	role := claims["role"].(string)
	return &model.TokenMetaData{
		UserID:  userID,
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

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(jwtKey), nil
}
