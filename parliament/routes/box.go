package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/controllers"
)

// Returns all routes for the ballot box model
func GetBoxRoutes(router *gin.RouterGroup) {
	boxRoutes := router.Group("/box")
	{
		// Routes for interacting with ballot boxes in the database
		boxRoutes.POST("/", controllers.CreateBox)
		boxRoutes.GET("/:id/", controllers.GetBoxById)
		boxRoutes.GET("/:id/:district/:number/", controllers.GetBoxByNumber)
		boxRoutes.PUT("/:id/", controllers.ChangeBox)
		boxRoutes.DELETE("/:id/:district/:number/", controllers.DeleteBox)
	}
}

// Returns all routes for the ballot box model
func GetBoxesRoutes(router *gin.RouterGroup) {
	boxRoutes := router.Group("/boxes")
	{
		// Routes for interacting with ballot boxes in the database
		boxRoutes.GET("/:city/", controllers.GetBoxesByCity)
		boxRoutes.GET("/:city/:district/", controllers.GetBoxesByDistrict)
		boxRoutes.GET("/:city/:district/:quarter/", controllers.GetBoxesByQuarter)
		// TODO: Get Boxes by Constituencies
	}
}
