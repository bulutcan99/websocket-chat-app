package utility

import (
	"fmt"
)

func VerifyRole(role string) (string, error) {
	switch role {
	case AdminRoleName:
	case UserRoleName:
	default:
		return "", fmt.Errorf("role '%v' does not exist", role)
	}

	return role, nil
}
