package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/presidency/controllers"
)

// Returns all routes for the quarter model
func GetQuarterRoutes(router *gin.RouterGroup) {
	quarterRoutes := router.Group("/quarter")
	{
		// Routes for interacting with quarters in the database
		quarterRoutes.POST("/", controllers.CreateQuarter)
		quarterRoutes.GET("/:id", controllers.GetQuarterById)
		quarterRoutes.GET("/:id/:district/:quarter/", controllers.GetQuarterByName)
		quarterRoutes.PUT("/:id/", controllers.ChangeQuarter)
		quarterRoutes.DELETE("/:id/", controllers.DeleteQuarter)
	}
}
