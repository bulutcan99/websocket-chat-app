package model

type Register struct {
	Name     string `json:"name" validate:"required,lte=25"`
	Surname  string `json:"surname" validate:"required,lte=25"`
	Nickname string `json:"nickname" validate:"required,lte=25"`
	Email    string `json:"email" validate:"required,nonzero,emailvalidator"`
	Password string `json:"password" validate:"required,passvalidator,lte=30"`
}

type SignIn struct {
	Email    string `json:"email" validate:"required,nonzero,emailvalidator"`
	Password string `json:"password" validate:"required,passvalidator,lte=30"`
}

type PasswordUpdate struct {
	NewPassword string `json:"new_password" validate:"required,passvalidator,lte=30"`
}
