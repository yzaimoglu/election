package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/controllers"
)

// Setup the user routes for the API
func GetUserRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/user")
	{
		// Routes for interacting with users in the database
		userRoutes.GET("/:id/", controllers.GetUser)
		userRoutes.PUT("/:id/email/", controllers.UpdateUserEmail)
		userRoutes.PUT("/:id/password/", controllers.UpdateUserPassword)
		userRoutes.PUT("/:id/affiliation/", controllers.UpdateUserAffiliation)
		userRoutes.PUT("/:id/role/", controllers.UpdateUserRole)
		userRoutes.PUT("/:id/lastseen/", controllers.UpdateUserLastseen)
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.DELETE("/:id/", controllers.DeleteUser)
	}
}
