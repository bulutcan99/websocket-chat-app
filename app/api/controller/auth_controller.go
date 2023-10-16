package controller

import (
	"context"
	"fmt"
	"github.com/bulutcan99/go-websocket/app/model"
	"github.com/bulutcan99/go-websocket/app/repository"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/helper"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type AuthController struct {
	repo *repository.AuthRepo
}

func NewAuthController(authRepo *repository.AuthRepo) *AuthController {
	return &AuthController{
		repo: authRepo,
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
		return fmt.Errorf("error while trying to set validation funcRole, %w", err)
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
