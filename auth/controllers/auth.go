package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp/totp"
	"github.com/yzaimoglu/election/auth/models"
	"github.com/yzaimoglu/election/auth/utilities"
	"gorm.io/gorm"
)

// Handle the Login
func LoginHandler(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Input required to login
	var input models.LoginInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check for an existing session
	if SessionExists(c) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "you are already logged in",
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

	// Object for the user to be logged in
	var user models.User

	// Check if user exists and return invalid username or email when not found
	if err := db.Where("username = ? AND email = ?", input.Username, input.Email).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  "no user with this username and email",
		})
		return
	}

	// Verify password
	if !utilities.CheckPassword(user.HashedPassword, input.PlainPassword) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  "wrong password",
		})
		return
	}

	// Try to verify totp
	if !totp.Validate(input.TOTP, string(utilities.FromBase64(user.TOTP))) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "could not verify the otp",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Generate Session
	session := createSessionObject(user.Id)

	// Create Session in Database
	if err := db.Create(&session).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Set Session Cookie
	c.SetCookie(
		"session-token",
		session.SessionToken,
		60*60*24, // in seconds (86400) --> one day
		"",
		utilities.GetEnv("AUTH_HOSTNAME", "localhost"),
		false,
		false,
	)

	// Return StatusOK
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "successfully logged in",
	})
}

// Handle the Login
func LogoutHandler(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	sessionToken, err := c.Cookie("session-token")
	if err != nil {
		sessionToken = c.GetHeader("Authentication-Session-Token")
	}

	// Create the session object
	var session models.Session

	// Check the database for the session token
	if err := db.Where("session_token = ?", sessionToken).First(&session).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   http.StatusNotFound,
			"message": "there is no such session",
		})
	}

	// Set Session Cookie for expiry
	c.SetCookie(
		"session-token",
		"expired",
		1,
		"/",
		utilities.GetEnv("AUTH_HOSTNAME", "localhost"),
		true,
		false,
	)

	// Delete Session from Database
	if err == nil {
		db.Delete(&session)
		// Return Success JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "successfully logged out",
			"status":  http.StatusOK,
		})
	}
}
