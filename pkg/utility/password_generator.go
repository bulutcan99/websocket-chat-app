package utility

import (
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func NormalizePassword(pass string) []byte {
	return []byte(pass)
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return false
	}
	if ok, _ := regexp.MatchString(`[a-z]`, password); !ok {
		return false
	}
	if ok, _ := regexp.MatchString(`[0-9]`, password); !ok {
		return false
	}

	return true
}

func GeneratePassword(pass string) (string, error) {
	err := isStrongPassword(pass)
	if err == false {
		return "", custom_error.PassError()
	}

	bytePwd := NormalizePassword(pass)
	hash, errCrypt := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if errCrypt != nil {
		return "", custom_error.ParseError()
	}

	return string(hash), nil
}

func ComparePasswords(hashedPwd, inputPwd string) bool {
	byteHash := NormalizePassword(hashedPwd)
	byteInput := NormalizePassword(inputPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
