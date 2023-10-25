package model

type Register struct {
	NameSurname string `json:"name_surname" validate:"required,lte=100"`
	Email       string `json:"email" validate:"required,nonzero,emailvalidator"`
	Password    string `json:"password" validate:"required,passvalidator,lte=30"`
	UserRole    string `json:"user_role" validate:"required,rolevalidator,lte=10"`
}

type SignIn struct {
	Email    string `json:"email" validate:"required,nonzero,emailvalidator"`
	Password string `json:"password" validate:"required,passvalidator,lte=30"`
}
