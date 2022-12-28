package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/presidency/middleware"
	"github.com/yzaimoglu/election/presidency/models"
	"github.com/yzaimoglu/election/presidency/routes"
	"github.com/yzaimoglu/election/presidency/utilities"
)

func main() {
	// Setup environment variables and some other things
	models.Setup()
	models.SetupCache()

	// Initialize the main router
	gin.SetMode(gin.ReleaseMode)
	mainRouter := gin.New()

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
		routes.GetCityRoutes(v1)
		routes.GetConstituencyRoutes(v1)
		routes.GetDistrictRoutes(v1)
		routes.GetQuarterRoutes(v1)
		routes.GetBoxRoutes(v1)
	}

	// Run server
	serverPort := fmt.Sprint(utilities.GetEnv("CB_PORT", fmt.Sprint(80)))
	if err := models.RedisSet("initialized_at", utilities.GetCurrentTime()); err != nil {
		fmt.Println("error initializing redis: " + err.Error())
	}
	fmt.Println("Cumhurbaşkanlığı server started running on port " + serverPort)
	mainRouter.Run(":" + serverPort)
}
