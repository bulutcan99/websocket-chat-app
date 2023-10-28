package utility

import (
	"errors"
	"go.uber.org/zap"
)

func VerifyRole(role string) (string, error) {
	switch role {
	case AdminRoleName:
		zap.S().Info("Admin user created successfully!")
	case UserRoleName:
		zap.S().Info("User created successfully!")
	default:
		return "", errors.New("role does not exist")
	}

	return role, nil
}
