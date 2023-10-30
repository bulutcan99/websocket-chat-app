package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id          int32      `json:"id" pg:",pk,type:serial"`
	Uuid        uuid.UUID  `json:"uuid" pg:"type:uuid,notnull,unique"`
	UserName    string     `json:"user_name" form:"user_name" binding:"required" pg:"type:varchar(25),notnull"`
	UserSurName string     `json:"user_surname" form:"user_surname" binding:"required" pg:"type:varchar(25),notnull"`
	Nickname    string     `json:"nickname" form:"nickname" binding:"required" pg:"type:varchar(25),notnull,unique"`
	Password    string     `json:"password" form:"password" binding:"required" pg:"type:varchar(255),notnull"`
	Email       string     `json:"email" form:"email" binding:"required" pg:"type:varchar(80),notnull,unique"`
	UserRole    string     `json:"user_role" form:"user_role" binding:"required" pg:"type:varchar(25),notnull"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type UserShown struct {
	Email       string    `json:"email" form:"email" binding:"required" pg:"type:varchar(80),notnull,unique"`
	UserName    string    `json:"user_name" form:"user_name" binding:"required" pg:"type:varchar(25),notnull"`
	UserSurName string    `json:"user_surname" form:"user_surname" binding:"required" pg:"type:varchar(25),notnull"`
	Nickname    string    `json:"nickname" form:"nickname" binding:"required" pg:"type:varchar(25),notnull,unique"`
	CreatedAt   time.Time `json:"created_at"`
	UserRole    string    `json:"user_role" validate:"required,rolevalidator,lte=10"`
	Status      int       `db:"user_status" json:"user_status" validate:"required,len=1"`
}
