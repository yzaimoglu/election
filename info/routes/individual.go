package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/info/controllers"
)

// Returns all routes for the individual model
func GetIndividualRoutes(router *gin.RouterGroup) {
	individualRoutes := router.Group("/individual")
	{
		// Routes for interacting with parties in the database
		individualRoutes.GET("/:id", controllers.GetIndividual)
		individualRoutes.POST("/", controllers.CreateIndividual)
		individualRoutes.PUT("/:id/", controllers.ChangeIndividual)
		individualRoutes.PUT("/:id/firstname/", controllers.ChangeIndividualFirstName)
		individualRoutes.PUT("/:id/lastname/", controllers.ChangeIndividualLastName)
		individualRoutes.PUT("/:id/birthdate/", controllers.ChangeIndividualBirthdate)
		individualRoutes.DELETE("/:id/", controllers.DeleteIndividual)
	}
}

// Returns all routes for the individuals model
func GetIndividualsRoutes(router *gin.RouterGroup) {
	individualsRoutes := router.Group("/individuals")
	{
		// Routes for interacting with parties in the database
		individualsRoutes.GET("/", controllers.GetIndividuals)
		individualsRoutes.GET("/:slice/", controllers.GetIndividualsBySlice)
	}
}
