package models

// Model for the totp verification
type TOTPVerification struct {
	Id       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Secret   string `json:"secret" validate:"required"`
	Image    string `json:"image" validate:"required"`
}

// Model for the totp verification input
type TOTPVerificationInput struct {
	Code string `json:"code" validate:"required"`
	TOTP string `json:"totp" validate:"required"`
}
