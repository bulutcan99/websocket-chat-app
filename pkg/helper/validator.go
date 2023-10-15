package helper

import (
	"fmt"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func EmailValidator(v interface{}, param string) error {
	value := v.(string)
	if strings.Contains(value, "@") && strings.Contains(value, ".") && len(value) <= 50 {
		parts := strings.Split(value, "@")
		if len(parts) == 2 {
			domainParts := strings.Split(parts[1], ".")
			if len(domainParts) >= 2 {
				return nil
			}
		}
	}

	return custom_error.ValidationError()
}

func RoleValidator(v interface{}, param string) error {
	role := v.(string)
	if role != "admin" && role != "user" {
		return custom_error.ValidationError()
	}
	return nil
}

func PasswordValidator(v interface{}, param string) error {
	password := v.(string)
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
