package utility

import (
	"fmt"
	"go.uber.org/zap"
)

func VerifyRole(role string) (string, error) {
	switch role {
	case AdminRoleName:
		zap.S().Infof("Admin user created successfully!")
	case UserRoleName:
		zap.S().Infof("User created successfully!")
	default:
		return "", fmt.Errorf("role '%v' does not exist", role)
	}

	return role, nil
}
