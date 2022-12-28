package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/controllers"
)

// Returns all routes for the district model
func GetDistrictRoutes(router *gin.RouterGroup) {
	districtRoutes := router.Group("/district")
	{
		// Routes for interacting with districts in the database
		districtRoutes.POST("/", controllers.CreateDistrict)
		districtRoutes.GET("/:id/", controllers.GetDistrictById)
		districtRoutes.GET("/:id/:district/", controllers.GetDistrictByName)
		districtRoutes.PUT("/:id/:district/", controllers.ChangeDistrict)
		districtRoutes.DELETE("/:id/", controllers.DeleteDistrict)
	}
}

// Returns all routes for the districts model
func GetDistrictsRoutes(router *gin.RouterGroup) {
	quartersRoutes := router.Group("/districts")
	{
		// Routes for interacting with districts in the database
		quartersRoutes.GET("/", controllers.GetDistricts)
		quartersRoutes.GET("/:city/", controllers.GetDistrictsByCity)
		quartersRoutes.GET("/:city/:constituency/", controllers.GetDistrictsByConstituency)
	}
}
