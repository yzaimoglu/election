package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/controllers"
)

// Returns all routes for the result model
func GetResultsRoutes(router *gin.RouterGroup) {
	resultsRoutes := router.Group("/results")
	{
		// Routes for interacting with results in the database
		resultsRoutes.GET("/:city/", controllers.GetResultsByCity)
		resultsRoutes.GET("/:city/:constituency/", controllers.GetResultsByConstituency)
		resultsRoutes.GET("/:city/:constituency/:district/", controllers.GetResultsByDistrict)
		resultsRoutes.GET("/:city/:constituency/:district/:quarter/", controllers.GetResultsByQuarter)
		resultsRoutes.GET("/:city/:constituency/:district/:quarter/:box/", controllers.GetResultsByBox)
	}
}
