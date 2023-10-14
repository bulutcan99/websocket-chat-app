package model

import "github.com/google/uuid"

type RenewToken struct {
	RefreshToken string `json:"refresh_token"`
}

type Error struct {
	Message string `json:"message"`
	Error   bool   `json:"error" default:"true"`
}

type Success struct {
	Message string `json:"message"`
	Error   bool   `json:"error" default:"false"`
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
}

type ChatRes struct {
	Message string `json:"message"`
	Error   bool   `json:"error" default:"false"`
	Chats   []Chat `json:"chats"`
}

type TokenMetaData struct {
	Name    string
	UserID  uuid.UUID
	Email   string
	Role    string
	Expires int64
}
