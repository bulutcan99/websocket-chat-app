package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"user_id" validate:"required,uuid"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Email        string    `json:"email" validate:"required,nonzero,emailvalidator"`
	NameSurname  string    `db:"name_surname" json:"name_surname" validate:"required,lte=100"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	Status       int       `db:"user_status" json:"user_status" validate:"required,len=1"`
	UserRole     string    `json:"user_role" validate:"required,rolevalidator,lte=10"`
}

type UserShown struct {
	Email       string    `json:"email" validate:"required,nonzero,emailvalidator"`
	NameSurname string    `db:"name_surname" json:"name_surname" validate:"required,lte=100"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UserRole    string    `json:"user_role" validate:"required,rolevalidator,lte=10"`
	Status      int       `db:"user_status" json:"user_status" validate:"required,len=1"`
}
