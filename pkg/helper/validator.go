package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *validator.Validate {
	valid := validator.New()

	_ = valid.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return false
		}
		return true
	})

	return valid
}

func ValidateErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
