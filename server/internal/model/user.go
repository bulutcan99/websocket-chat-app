package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id           int32      `json:"id" pg:",pk,type:serial"`
	UUID         uuid.UUID  `json:"uuid" pg:"type:uuid,notnull,unique"`
	UserName     string     `json:"user_name" form:"user_name" binding:"required" pg:"type:varchar(25),notnull"`
	UserSurName  string     `json:"user_surname" form:"user_surname" binding:"required" pg:"type:varchar(25),notnull"`
	Nickname     string     `json:"nickname" form:"nickname" binding:"required" pg:"type:varchar(25),notnull,unique"`
	Passwordhash string     `json:"passwordhash" form:"passwordhash" binding:"required" pg:"type:varchar(255),notnull"`
	Email        string     `json:"email" form:"email" binding:"required" pg:"type:varchar(80),notnull,unique"`
	UserRole     string     `json:"user_role" form:"user_role" binding:"required" pg:"type:varchar(25),notnull"`
	Status       int        `db:"user_status" json:"user_status" validate:"required,len=1"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	BlockedAt    *time.Time `json:"blocked_at"`
}

type UserShown struct {
	Id          int32     `json:"id"`
	Email       string    `json:"email"`
	UserName    string    `json:"user_name"`
	UserSurName string    `json:"user_surname"`
	Nickname    string    `json:"nickname"`
	CreatedAt   time.Time `json:"created_at"`
	UserRole    string    `json:"user_role"`
	Status      int       `db:"user_status"`
}
