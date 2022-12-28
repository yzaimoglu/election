package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/presidency/controllers"
)

// Returns all routes for the district model
func GetDistrictRoutes(router *gin.RouterGroup) {
	districtRoutes := router.Group("/district")
	{
		// Routes for interacting with districts in the database
		districtRoutes.POST("/", controllers.CreateDistrict)
		districtRoutes.GET("/:id/", controllers.GetDistrictById)
		districtRoutes.GET("/:id/:district/", controllers.GetDistrictByName)
		districtRoutes.PUT("/:id/", controllers.ChangeDistrict)
		districtRoutes.DELETE("/:id/", controllers.DeleteDistrict)
	}
}
