package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/presidency/controllers"
)

// Returns all routes for the city model
func GetCityRoutes(router *gin.RouterGroup) {
	cityRoutes := router.Group("/city")
	{
		// Routes for interacting with cities in the database
		cityRoutes.POST("/", controllers.CreateCity)
		cityRoutes.GET("/:id", controllers.GetCity)
		cityRoutes.PUT("/:id/", controllers.ChangeCity)
		cityRoutes.DELETE("/:id/", controllers.DeleteCity)
	}
}
