package helper

import (
	"errors"
	"fmt"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"time"
)

func EmailValidator(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email address")
	}
	fmt.Println("AS")
	return nil
}

func RoleValidator(role string) error {
	if role != "admin" && role != "user" {
		return custom_error.ValidationError()
	}
	return nil
}

func PasswordValidator(password string) error {
	if len(password) < 9 {
		return fmt.Errorf("Password must be at least 9 characters long")
	}
	// Check for uppercase letter (A-Z)
	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return fmt.Errorf("Password must contain at least one uppercase letter")
	}

	// Check for lowercase letter (a-z)
	if ok, _ := regexp.MatchString(`[a-z]`, password); !ok {
		return fmt.Errorf("Password must contain at least one lowercase letter")
	}

	// Check for digit (0-9)
	if ok, _ := regexp.MatchString(`[0-9]`, password); !ok {
		return fmt.Errorf("Password must contain at least one digit")
	}

	return nil // Password is valid
}

// TokenValidator checks if token is valid UUID
func TokenValidator(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return custom_error.ValidationError()
	}
	_, err := uuid.Parse(v.(string))
	if err != nil {
		return custom_error.ValidationError()
	}
	return nil
}

// DateValidator checks if date provided is not less than today's date
func DateValidator(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if date, ok := st.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return custom_error.ValidationError()
		}
	} else {
		return custom_error.ValidationError()
	}
	return nil
}
