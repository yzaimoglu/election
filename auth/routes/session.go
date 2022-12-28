package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/controllers"
)

// Setup the session routes for the API
func GetSessionRoutes(router *gin.RouterGroup) {
	sessionRoutes := router.Group("/")
	{
		// Routes for interacting with the sessions in the database
		sessionRoutes.GET("/session/:session/", controllers.GetSessionById)
		sessionRoutes.GET("/session/", controllers.GetSessionByCookie)
		sessionRoutes.GET("/sessions/:id/", controllers.GetSessionsOfUser)
	}
}
