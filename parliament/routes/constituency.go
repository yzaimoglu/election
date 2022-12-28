package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/controllers"
)

// Returns all routes for the constituency model
func GetConstituencyRoutes(router *gin.RouterGroup) {
	constituencyRoutes := router.Group("/constituency")
	{
		// Routes for interacting with constituencies in the database
		constituencyRoutes.POST("/", controllers.CreateConstituency)
		constituencyRoutes.GET("/:id/", controllers.GetConstituency)
		constituencyRoutes.PUT("/:id/", controllers.ChangeConstituency)
		constituencyRoutes.DELETE("/:id/", controllers.DeleteConstituency)
	}
}

// Returns all routes for the constituency model
func GetConstituenciesRoutes(router *gin.RouterGroup) {
	constituenciesRoutes := router.Group("/constituencies")
	{
		// Routes for interacting with constituencies in the database
		constituenciesRoutes.GET("/", controllers.GetConstituencies)
		constituenciesRoutes.GET("/:city/", controllers.GetConstituenciesByCity)
	}
}
