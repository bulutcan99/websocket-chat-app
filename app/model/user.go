package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Email        string    `db:"email" json:"email" validate:"required,email,lte=255"`
	Name         string    `db:"name" json:"name" validate:"required,lte=50"`
	Surname      string    `db:"surname" json:"surname" validate:"required,lte=50"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	Password     string    `json:"password,omitempty" validate:"lte=30"`
	Status       int       `db:"user_status" json:"user_status" validate:"required,len=1"`
	Role         string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
}

type UserShown struct {
	Email     string    `db:"email" json:"email" validate:"required,email,lte=255"`
	Name      string    `db:"name" json:"name" validate:"required,lte=50"`
	Surname   string    `db:"surname" json:"surname" validate:"required,lte=50"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserRole  string    `db:"user_role" json:"user_role" validate:"required,lte=25"`
}
