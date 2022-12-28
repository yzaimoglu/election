package models

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/yzaimoglu/election/updater/utilities"
)

// Setup Environment variables and make sure database is initialized
func Setup() {
	// Load the environment variables
	godotenv.Load()

	// Check if system is in Debug Mode
	DEBUG := utilities.GetEnv("UPDATER_DEBUG", "false")

	// Sleep to make sure that the Database is initialized beforehand
	if DEBUG != "false" {
		time.Sleep(20 * time.Second)
	}
}
