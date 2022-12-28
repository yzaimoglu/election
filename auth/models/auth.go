package models

// Input required to login
type LoginInput struct {
	Username      string `json:"username" validate:"required"`
	Email         string `json:"email" validate:"required"`
	PlainPassword string `json:"plainpassword" validate:"required"`
	TOTP          string `json:"totp" validate:"required"`
}

// Model for the session object
type Session struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"userid"`
	SessionToken string `json:"sessiontoken"`
	CreatedAt    int64  `json:"createdat"`
	ExpiresAt    int64  `json:"expiresat"`
}
