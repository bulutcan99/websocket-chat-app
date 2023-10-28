package model

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
