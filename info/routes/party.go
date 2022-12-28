package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/info/controllers"
)

// Returns all routes for the party model
func GetPartyRoutes(router *gin.RouterGroup) {
	partyRoutes := router.Group("/party")
	{
		// Routes for interacting with parties in the database
		partyRoutes.GET("/:id", controllers.GetParty)
		partyRoutes.POST("/", controllers.CreateParty)
		partyRoutes.PUT("/:id/", controllers.ChangeParty)
		partyRoutes.PUT("/:id/name/", controllers.ChangePartyName)
		partyRoutes.PUT("/:id/abbreviation/", controllers.ChangePartyAbbreviation)
		partyRoutes.PUT("/:id/leader/", controllers.ChangePartyLeader)
		partyRoutes.PUT("/:id/logo/", controllers.ChangePartyLogo)
		partyRoutes.PUT("/:id/color/", controllers.ChangePartyColor)
		partyRoutes.DELETE("/:id/", controllers.DeleteParty)
	}
}

// Returns all routes for the parties model
func GetPartiesRoutes(router *gin.RouterGroup) {
	partiesRoutes := router.Group("/parties")
	{
		// Routes for interacting with parties in the database
		partiesRoutes.GET("/", controllers.GetParties)
		partiesRoutes.GET("/:slice/", controllers.GetPartiesBySlice)
	}
}
