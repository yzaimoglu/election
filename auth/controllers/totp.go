package controllers

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/yzaimoglu/election/auth/models"
	"github.com/yzaimoglu/election/auth/utilities"
	"gorm.io/gorm"
)

// Verify the totp
func VerifyTOTP(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)

	// Create the verificationcode model
	var verification models.TOTPVerification
	var verificationInput models.TOTPVerificationInput

	// Bind the input from the request body to the Input object
	if err := c.ShouldBindJSON(&verificationInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(verificationInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Check the database for the verification code
	if err := db.Where("code = ?", verificationInput.Code).First(&verification).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "could not find such a code in the database",
			"status":  http.StatusNotFound,
		})
		return
	}

	// Try to verify totp
	if !totp.Validate(verificationInput.TOTP, string(utilities.FromBase64(verification.Secret))) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "could not verify the totp",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// User to be updated
	var user models.User

	// Find the User and return 500 when not found
	if err := db.Where("username = ?", verification.Username).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Create an updated user object
	var updatedUser models.User = user
	updatedUser.TOTP = verification.Secret

	// Save the Updated User to the Database
	db.Model(&user).Updates(updatedUser)

	// Delete the verification code from the database
	if err := db.Where("code = ?", verification.Code).Delete(&verification); err.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  err.Error,
		})
		return
	}

	// Return verified totp
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "user totp has been successfully verified",
	})
}

// Get the image
func GetImage(c *gin.Context) {
	// Get the database connection from the context
	db := c.MustGet("db").(*gorm.DB)
	// Get the verification code from the url
	verificationCode := c.Param("verification")

	// Create the verificationcode model
	var verification models.TOTPVerification

	// Check the database for the verification code
	if err := db.Where("code = ?", verificationCode).First(&verification).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "there is no such image",
			"status":  http.StatusNotFound,
		})
		return
	}

	// Create the file
	file, err := os.Create(verificationCode + ".png")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}
	// Write to the file
	if _, err := file.Write(utilities.FromBase64(verification.Image)); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}
	file.Close()

	// Return the file
	c.File(verificationCode + ".png")

	// Remove the file
	if err := os.Remove(verificationCode + ".png"); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}
}

// Create TOTP
func CreateTOTP(accountemail string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "secim2023.org",
		AccountName: accountemail,
	})
}

// Get Image bytes
func GetImageBytes(key *otp.Key) ([]byte, error) {
	var buf bytes.Buffer
	img, err := key.Image(500, 500)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)
	return buf.Bytes(), nil
}

// Get Image bytes
func GetImageStruct(key *otp.Key) (image.Image, error) {
	img, err := key.Image(500, 500)
	if err != nil {
		return nil, err
	}
	return img, nil
}
