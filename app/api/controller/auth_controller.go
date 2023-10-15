package controller

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/app/model"
	"github.com/bulutcan99/go-websocket/pkg/db/sql"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/bulutcan99/go-websocket/pkg/utility"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func UserSignUp(c *fiber.Ctx) error {
	signUp := &model.Register{}
	if err := c.BodyParser(signUp); err != nil {
		return custom_error.ParseError()
	}

	// validator.SetValidationFunc("emailvalidator", helper.EmailValidator)
	// validator.SetValidationFunc("passvalidator", helper.PasswordValidator)
	// validator.SetValidationFunc("rolevalidator", helper.RoleValidator)
	// if err := validator.Validate(signUp); err == nil {
	// 	fmt.Println("Validated.")
	// } else {
	// 	return custom_error.ValidationError()
	// }

	db := sql.SqlQueryInjection()
	// errDb := db.UserQueries.DB.Ping()
	// if errDb != nil {
	// 	return custom_error.ConnectionError()
	// }

	fmt.Println("Queries injected")
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
	if errCreate := db.AuthQueries.CreateUser(user); errCreate != nil {
		return custom_error.DatabaseError()
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"user":  user,
	})
}
