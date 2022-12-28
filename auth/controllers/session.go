package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/models"
	"github.com/yzaimoglu/election/auth/utilities"
	"gorm.io/gorm"
)

// Get the active sessions of the user
func GetSessionsOfUser(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Create the user object
	var sessions []models.Session

	// Find all sessions
	if err := db.Where("user_id = ?", id).Find(&sessions).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "no active sessions found",
			"status":  http.StatusNotFound,
		})
		return
	}

	// Return userid and sessions
	c.JSON(http.StatusOK, gin.H{
		"userid":   id,
		"sessions": sessions,
	})
}

// Return Session with user of provided session id or token
func GetSessionById(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	sessionToken := c.Param("session")

	// Create the session object
	var session models.Session

	// Check the database for the session token
	if err := db.Where("id = ? OR session_token = ?", sessionToken, sessionToken).First(&session).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "you are not logged in",
			"status":  http.StatusNotFound,
		})
		return
	}

	// Check if session has expired if so delete session
	if utilities.GetCurrentTime() >= session.ExpiresAt {
		var message string
		var status int
		if err := db.Delete(&session); err.Error != nil {
			message = "internal server error"
			status = http.StatusInternalServerError
		} else {
			message = "the session has expired"
			status = http.StatusNotFound
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": message,
			"status":  status,
		})
		return
	}

	// Initialize the User object
	var user models.User

	// Find the User and return 404 including error when not found
	if err :=
		db.Where("id = ?", session.UserId).
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

	// Return session and user
	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"user":    userInformation,
	})
}

// Return Session with user
func GetSessionByCookie(c *gin.Context) {
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
			"message": "you are not logged in",
			"status":  http.StatusNotFound,
		})
		return
	}

	// Check if Tokens are identical
	if sessionToken != session.SessionToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Check if session has expired if so delete session
	if utilities.GetCurrentTime() >= session.ExpiresAt {
		var message string
		var status int
		if err := db.Delete(&session); err.Error != nil {
			message = "internal server error"
			status = http.StatusInternalServerError
		} else {
			message = "you are not logged in"
			status = http.StatusNotFound
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": message,
			"status":  status,
		})
		return
	}

	// Initialize the User object
	var user models.User

	// Find the User and return 404 including error when not found
	if err :=
		db.Where("id = ?", session.UserId).
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

	// Return session and user
	c.JSON(http.StatusOK, gin.H{
		"session": session,
		"user":    userInformation,
	})
}

// Check if a session exists
func SessionExists(c *gin.Context) bool {
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
		return false
	}

	// Check if Tokens are identical
	if sessionToken != session.SessionToken {
		return false
	}

	// Check if session has expired if so delete session
	if utilities.GetCurrentTime() >= session.ExpiresAt {
		db.Delete(&session)
		return false
	}

	return true
}

// Create a session object
func createSessionObject(userId int64) models.Session {
	sessionToken, _ := utilities.GenerateRandomBytes(256)
	var session models.Session
	session.UserId = userId
	session.CreatedAt = utilities.GetCurrentTime()
	session.ExpiresAt = session.CreatedAt + (60 * 60 * 24 * 1000)
	session.SessionToken = utilities.ToBase64(sessionToken)
	return session
}
