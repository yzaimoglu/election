package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/presidency/controllers"
)

// Returns all routes for the ballot box model
func GetBoxRoutes(router *gin.RouterGroup) {
	boxRoutes := router.Group("/box")
	{
		// Routes for interacting with ballot boxes in the database
		boxRoutes.POST("/", controllers.CreateBox)
		boxRoutes.GET("/:id", controllers.GetBoxById)
		boxRoutes.GET("/:id/:district/:number/", controllers.GetBoxByNumber)
		boxRoutes.PUT("/:id/", controllers.ChangeBox)
		boxRoutes.DELETE("/:id/:district/:number/", controllers.DeleteBox)
	}
}
