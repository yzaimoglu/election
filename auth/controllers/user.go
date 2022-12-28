package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	mail "github.com/xhit/go-simple-mail/v2"
	"github.com/yzaimoglu/election/auth/models"
	"github.com/yzaimoglu/election/auth/utilities"
	"gorm.io/gorm"
)

// Check if the User exists
func UserExists(c *gin.Context, id string) bool {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Initialize the User object
	var user models.User

	// Find the User and return true or false
	if err :=
		db.Where("id = ? OR username = ? OR email = ?", id, id, id).
			First(&user).Error; err != nil {
		return false
	}
	return true
}

// Check if the Username exists
func UsernameExists(c *gin.Context, username string) bool {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Initialize the User object
	var user models.User

	// Find the User and return true or false
	if err :=
		db.Where("username = ?", username).
			First(&user).Error; err != nil {
		return false
	}
	return true
}

// Check if the Email exists
func EmailExists(c *gin.Context, email string) bool {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Initialize the User object
	var user models.User

	// Find the User and return true or false
	if err :=
		db.Where("email = ?", email).
			First(&user).Error; err != nil {
		return false
	}
	return true
}

// Find user
func GetUser(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 including error when not found
	if err :=
		db.Where("id = ? OR username = ? OR email = ?", id, id, id).
			First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// User without sensitive information
	userInformation := models.UserInformation{
		Id:          user.Id,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		LastSeen:    user.LastSeen,
		Role:        user.Role,
		Affiliation: user.Affiliation,
	}

	// Return the User
	c.JSON(http.StatusOK, userInformation)
}

// Create a new user
func CreateUser(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Initialize the CreateAnimecloUser Model
	var createUser models.CreateUser

	// Get the input from the Request Body and Decode into CreateUser Object
	if err := json.NewDecoder(c.Request.Body).Decode(&createUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(createUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check if there is a user with the specified email
	if EmailExists(c, createUser.Email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "email already exists",
		})
		return
	}

	// Create username and get rid of turkish characters from username
	username := strings.ToLower(strings.ReplaceAll(createUser.FirstName, " ", "-") + "." + strings.ReplaceAll(createUser.LastName, " ", "-"))
	username = strings.ReplaceAll(username, "ş", "s")
	username = strings.ReplaceAll(username, "ğ", "g")
	username = strings.ReplaceAll(username, "ç", "c")
	username = strings.ReplaceAll(username, "ı", "i")
	username = strings.ReplaceAll(username, "ö", "o")
	username = strings.ReplaceAll(username, "ü", "u")

	i := 1
	usernameWithoutNumbers := username

	// Check if username already exists, if it does append append an integer
	for UsernameExists(c, username) {
		username = usernameWithoutNumbers + fmt.Sprint(i)
		i++
	}

	// Create TOTP
	totp, err := CreateTOTP(createUser.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Create the new user model
	user := models.User{
		Username:       username,
		FirstName:      createUser.FirstName,
		LastName:       createUser.LastName,
		Email:          createUser.Email,
		Role:           "Kullanıcı",
		Affiliation:    createUser.Affiliation,
		CreatedAt:      utilities.GetCurrentTime(),
		LastSeen:       -1,
		HashedPassword: utilities.HashPassword(createUser.PlainPassword),
		TOTP:           utilities.ToBase64([]byte(totp.Secret())),
	}

	// Create the totp verification object
	imageBytes, err := GetImageBytes(totp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	bytes, err := utilities.GenerateRandomBytes(50)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	totpVerification := models.TOTPVerification{
		Username: username,
		Code:     utilities.ToBase64(bytes) + uuid.New().String(),
		Secret:   utilities.ToBase64([]byte(totp.Secret())),
		Image:    utilities.ToBase64(imageBytes),
	}

	// Create Mail info
	// TODO: Will be done in the emailer later on
	createMail := models.Mail{
		From:    "election/auth Tracker - Authentication <user@user>",
		To:      createUser.Email,
		Subject: "election/auth Tracker Authentication Requirement",
		Body:    "https:///localhost/totp/" + totpVerification.Code,
		Credentials: models.MailCredentials{
			Username: utilities.GetEnv("SMTP_USER", "user@user"),
			Password: utilities.GetEnv("SMTP_PASSWORD", "password"),
		},
		Server: models.MailServer{
			Host:       utilities.GetEnv("SMTP_HOST", "smtp_host"),
			Port:       587,
			Encryption: mail.EncryptionSTARTTLS,
		},
	}

	// Specify mailserver options
	mailServer := mail.NewSMTPClient()
	mailServer.Host = createMail.Server.Host
	mailServer.Port = createMail.Server.Port
	mailServer.Username = createMail.Credentials.Username
	mailServer.Password = createMail.Credentials.Password
	mailServer.Encryption = createMail.Server.Encryption

	// Connect to mailserver
	smtpClient, err := mailServer.Connect()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Set email info
	email := mail.NewMSG()
	email.SetFrom(createMail.From)
	email.AddTo(createMail.To)
	email.SetSubject(createMail.Subject)
	email.SetBody(mail.TextHTML, createMail.Body)

	// Send mail
	if err := email.Send(smtpClient); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Create the TOTP Verification in the Database
	db.Select("Id", "Username", "Code", "Secret", "Image").Create(&totpVerification)

	// Create the User in the Database
	db.Select("Id", "Username", "FirstName", "LastName", "Email", "HashedPassword", "CreatedAt", "LastSeen",
		"Role", "Affiliation", "TOTP").Create(&user)

	// Return the newly created user
	c.JSON(http.StatusOK, user)
}

// Update the email of a user
func UpdateUserEmail(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Initialize the email input
	var input models.UpdateEmailInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check for duplicate email
	if EmailExists(c, input.Email) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusBadRequest,
			"error":  "email already in use",
		})
		return
	}

	// Check if input is the same as the email of the user
	if user.Email == input.Email {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusBadRequest,
			"error":  "cannot change to same email",
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.Email = input.Email

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Return the updated fields of the user object
	c.JSON(http.StatusOK, input)
}

// Update the email of a user
func UpdateUserPassword(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Initialize the password input
	var input models.UpdatePasswordInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Verify old password
	if !utilities.CheckPassword(user.HashedPassword, input.OldPassword) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusUnauthorized,
			"error":  "unauthorized",
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check if input is the same as the email of the user
	if input.OldPassword == input.NewPassword {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusBadRequest,
			"error":  "cannot change to same password",
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.HashedPassword = utilities.HashPassword(input.NewPassword)

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Return the updated fields of the user object
	c.JSON(http.StatusOK, input)
}

// Update the affiliation of a user
func UpdateUserAffiliation(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Initialize the affiliation input
	var input models.UpdateAffiliationInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check if input is the same as the affiliation of the user
	if user.Affiliation == input.Affiliation {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusBadRequest,
			"error":  "cannot change to same affiliation",
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.Affiliation = input.Affiliation

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Return the updated fields of the user object
	c.JSON(http.StatusOK, input)
}

// Update the role of a user
func UpdateUserRole(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Initialize the role input
	var input models.UpdateRoleInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check if input is the same as the role of the user
	if user.Role == input.Role {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusBadRequest,
			"error":  "cannot change to same role",
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.Role = input.Role

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Return the updated fields of the user object
	c.JSON(http.StatusOK, input)
}

// Update the last seen time of a user
func UpdateUserLastseen(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Initialize the lastseen input
	var input models.UpdateLastseenInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.LastSeen = input.LastSeen

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Return the updated fields of the user object
	c.JSON(http.StatusOK, input)
}

// Deletes the user
func DeleteUser(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Initialize the User object
	var user models.User

	// Find the User and return 404 when not found
	if err := db.Where("id = ? OR username = ? OR email = ?", id, id, id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error(),
		})
		return
	}

	// Delete the user and return
	if err := db.Delete(&user); err.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error,
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
