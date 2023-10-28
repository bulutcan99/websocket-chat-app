package utility

import (
	"golang.org/x/crypto/bcrypt"
)

func NormalizePassword(pass string) []byte {
	return []byte(pass)
}

func GeneratePassword(pass string) string {
	bytePwd := NormalizePassword(pass)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return err.Error()
	}

	return string(hash)
}

func ComparePasswords(hashedPwd, inputPwd string) bool {
	byteHash := NormalizePassword(hashedPwd)
	byteInput := NormalizePassword(inputPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
