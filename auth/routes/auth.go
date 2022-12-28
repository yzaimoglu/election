package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/controllers"
)

// Setup the auth routes for the API
func GetAuthRoutes(router *gin.RouterGroup) {
	authRoutes := router.Group("/")
	{
		// Routes for authentication procedures
		authRoutes.POST("/login/", controllers.LoginHandler)
		authRoutes.DELETE("/logout/", controllers.LogoutHandler)
	}
}
