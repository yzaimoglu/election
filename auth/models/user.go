package models

// Model for the user object
type User struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedpassword"`
	CreatedAt      int64  `json:"createdat"`
	LastSeen       int64  `json:"lastseen"`
	Role           string `json:"role"`
	Affiliation    string `json:"affiliation"`
	TOTP           string `json:"totp"`
}

// Model for user information (User object without sensitive information)
type UserInformation struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdat"`
	LastSeen    int64  `json:"lastseen"`
	Role        string `json:"role"`
	Affiliation string `json:"affiliation"`
}

// Model to use for user creation
type CreateUser struct {
	FirstName     string `json:"firstname" validate:"required"`
	LastName      string `json:"lastname" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	PlainPassword string `json:"plainpassword" validate:"required"`
	Affiliation   string `json:"affiliation" validate:"required"`
}

// Model for the update email input
type UpdateEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

// Model for the update password input
type UpdatePasswordInput struct {
	OldPassword string `json:"oldpassword" validate:"required"`
	NewPassword string `json:"newpassword" validate:"required"`
}

// Model for the update affiliation input
type UpdateAffiliationInput struct {
	Affiliation string `json:"affiliation" validate:"required"`
}

// Model for the update role input
type UpdateRoleInput struct {
	Role string `json:"role" validate:"required"`
}

// Model for the update lastseen input
type UpdateLastseenInput struct {
	LastSeen int64 `json:"lastseen" validate:"required,numeric"`
}

// Model for the update totp input
type UpdateTOTPInput struct {
	TOTP string `json:"totp" validate:"required"`
}
