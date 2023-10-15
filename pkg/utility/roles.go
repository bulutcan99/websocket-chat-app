package utility

import (
	"fmt"
)

func VerifyRole(role string) (string, error) {
	switch role {
	case AdminRoleName:
		fmt.Println("Admin user created successfully!")
	case UserRoleName:
		fmt.Println("User created successfully!")
	default:
		return "", fmt.Errorf("role '%v' does not exist", role)
	}

	return role, nil
}
