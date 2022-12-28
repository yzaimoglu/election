package models

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/yzaimoglu/election/auth/utilities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Setup the database connection with the environment variables
func SetupDB() *gorm.DB {
	// Load the environment variables
	godotenv.Load()

	// Check if system is in Debug Mode
	DEBUG := utilities.GetEnv("AUTH_DEBUG", "false")

	// Sleep to make sure that the Database is initialized beforehand
	if DEBUG != "false" {
		time.Sleep(20 * time.Second)
	}

	// Database credentials from the environment variables
	USER := utilities.GetEnv("AUTH_DB_USER", "root")
	PASS := utilities.GetEnv("AUTH_DB_PASSWORD", "password")
	HOST := utilities.GetEnv("AUTH_DB_HOST", "localhost")
	PORT := utilities.GetEnv("AUTH_DB_PORT", "3306")
	DBNAME := utilities.GetEnv("AUTH_DB_DATABASE", "authentication")

	// Establishing the connection to the database
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open(mysql.Open(URL))

	// Check for errors
	if err != nil {
		panic("failed to connect database")
		// panic(err)
	}

	return db
}
