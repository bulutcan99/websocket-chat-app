package custom_error

import (
	"fmt"
)

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func NewCustomError(code int, message string) error {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func ReqBodyDataError() error {
	return NewCustomError(400, "Request body data is not valid!")
}

func DatabaseError() error {
	return NewCustomError(400, "Database error!")
}

func ValidationError() error {
	return NewCustomError(400, "There is an error in validation.")
}

func PassError() error {
	return NewCustomError(400, "Your password is invalid!")
}

func ParseError() error {
	return NewCustomError(400, "Parsing error")
}

func ConnectionError() error {
	return NewCustomError(400, "Oops... Connection failed!")
}
