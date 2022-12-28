package models

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/yzaimoglu/election/info/utilities"
)

// Setup Environment variables and make sure database is initialized
func Setup() {
	// Load the environment variables
	godotenv.Load()

	// Check if system is in Debug Mode
	DEBUG := utilities.GetEnv("BILGI_DEBUG", "false")

	// Sleep to make sure that the Database is initialized beforehand
	if DEBUG != "false" {
		time.Sleep(20 * time.Second)
	}
}

// Model for the color used in the other models
type Color struct {
	Hex string `json:"hex" bson:"hex"`
	JS  string `json:"js" bson:"js"`
}
