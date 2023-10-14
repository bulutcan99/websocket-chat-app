package model

type Register struct {
	Email    string `json:"email" validate:"required,email,lte=50"`
	Password string `json:"password" validate:"required,lte=30"`
	UserRole string `json:"user_role" validate:"required,lte=25"`
}

type SignIn struct {
	Email    string `json:"email" validate:"required,email,lte=50"`
	Password string `json:"password" validate:"required,lte=30"`
}
