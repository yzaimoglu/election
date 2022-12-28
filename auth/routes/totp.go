package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/auth/controllers"
)

// Setup the totp routes for the API
func GetTOTPRoutes(router *gin.RouterGroup) {
	totpRoutes := router.Group("/totp")
	{
		// Routes for interacting with the totp verification codes in the database
		totpRoutes.POST("/", controllers.VerifyTOTP)
		totpRoutes.GET("/:verification/", controllers.GetImage)
	}
}
