package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/middleware"
	"github.com/yzaimoglu/election/auth/models"
	"github.com/yzaimoglu/election/auth/routes"
	"github.com/yzaimoglu/election/auth/utilities"
)

func main() {
	// Setup the Database and create the tables
	db := models.SetupDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Session{})
	db.AutoMigrate(&models.TOTPVerification{})

	// Initialize the main router
	gin.SetMode(gin.ReleaseMode)
	mainRouter := gin.New()

	// Add the database to the context
	mainRouter.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// Setup default security measures
	mainRouter.Use(middleware.Default())

	// Setup the CORS middleware
	mainRouter.Use(middleware.CORSMiddleware)

	// Setup the Basic and Security Middleware provided by Gin
	mainRouter.Use(gin.Logger())
	mainRouter.Use(gin.Recovery())

	// Standard NoRoute Response
	mainRouter.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "not found"})
	})

	// Standard NoMethod Response
	mainRouter.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"status": http.StatusMethodNotAllowed, "error": "method not allowed"})
	})

	// Set the Favicon
	mainRouter.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// Create the main Route group for the API
	v1 := mainRouter.Group("/v1")
	{
		routes.GetUserRoutes(v1)
		routes.GetAuthRoutes(v1)
		routes.GetSessionRoutes(v1)
		routes.GetTOTPRoutes(v1)
	}

	// Run server
	serverPort := fmt.Sprint(utilities.GetEnv("AUTH_PORT", fmt.Sprint(80)))
	fmt.Println("Authentication server started running on port " + serverPort)
	mainRouter.Run(":" + serverPort)
}
