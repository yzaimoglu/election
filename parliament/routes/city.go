package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/controllers"
)

// Returns all routes for the city model
func GetCityRoutes(router *gin.RouterGroup) {
	cityRoutes := router.Group("/city")
	{
		// Routes for interacting with cities in the database
		cityRoutes.POST("/", controllers.CreateCity)
		cityRoutes.GET("/:id/", controllers.GetCity)
		cityRoutes.PUT("/:id/", controllers.ChangeCity)
		cityRoutes.DELETE("/:id/", controllers.DeleteCity)
	}
}

// Returns all routes for the cities model
func GetCitiesRoutes(router *gin.RouterGroup) {
	citiesRoutes := router.Group("/cities")
	{
		// Routes for interacting with cities in the database
		citiesRoutes.GET("/", controllers.GetCities)
	}
}
