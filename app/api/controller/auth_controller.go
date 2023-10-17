package controller

import (
	"context"
	"fmt"
	"github.com/bulutcan99/go-websocket/app/model"
	"github.com/bulutcan99/go-websocket/app/repository"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/helper"
	"github.com/bulutcan99/go-websocket/pkg/jwt"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthInterface interface {
	UserSignUp(c *fiber.Ctx) error
	UserSignIn(c *fiber.Ctx) error
}

type AuthController struct {
	repo       *repository.AuthRepo
	redisCache *redis.Client
}

func NewAuthController(authRepo *repository.AuthRepo, redisC *redis.Client) *AuthController {
	return &AuthController{
		repo:       authRepo,
		redisCache: redisC,
	}
}

func (ac *AuthController) UserSignUp(c *fiber.Ctx) error {
	signUp := &model.Register{}
	if err := c.BodyParser(signUp); err != nil {
		return custom_error.ParseError()
	}

	signUpCtx := context.Background()
	err := helper.EmailValidator(signUp.Email)
	if err != nil {
		return fmt.Errorf("error while trying to set validation funcEmail, %w", err)
	}
	err = helper.PasswordValidator(signUp.Password)
	if err != nil {
		return fmt.Errorf("error while trying to set validation funcPass, %w", err)
	}
	err = helper.RoleValidator(signUp.UserRole)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Role must be admin or user",
		})
	}

	_, errVerify := utility.VerifyRole(signUp.UserRole)
	if errVerify != nil {
		return custom_error.ValidationError()
	}

	user := model.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		NameSurname:  signUp.NameSurname,
		Email:        signUp.Email,
		PasswordHash: utility.GeneratePassword(signUp.Password),
		Status:       1,
		UserRole:     signUp.UserRole,
	}

	if errCreate := ac.repo.CreateUser(signUpCtx, user); errCreate != nil {
		return custom_error.DatabaseError()
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}

func (ac *AuthController) UserSignIn(c *fiber.Ctx) error {
	signIn := &model.SignIn{}
	if err := c.BodyParser(signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "error while trying to parse body",
		})
	}

	getUser, err := ac.repo.GetUserSignByEmail(context.Background(), signIn.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Given email is not found",
		})
	}

	isComparedUserPass := utility.ComparePasswords(getUser.PasswordHash, signIn.Password)
	if !isComparedUserPass {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Given password is not correct",
		})
	}

	role := getUser.UserRole

	accessToken, err := jwt.GenerateNewTokens(getUser.ID.String(), role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to generate new tokens",
		})
	}

	userId := getUser.ID.String()
	errSaveToRedis := ac.redisCache.Set(context.Background(), userId, accessToken.Refresh, 0).Err()
	if errSaveToRedis != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "There is a problem while trying to save redis",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Logged In Successfully!",
		"tokens": fiber.Map{
			"access":  accessToken.Access,
			"refresh": accessToken.Refresh,
		},
	})
}
