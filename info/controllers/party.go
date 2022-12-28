package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/yzaimoglu/election/info/models"
	"github.com/yzaimoglu/election/info/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Returns a single party with a specific id
func GetParty(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	var party models.Party

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Find the party in the database
	result := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").FindOne(ctx, &bson.M{"_id": objId})

	// Check if there is a party with that id
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no party with that id found",
		})
		return
	}

	// Decode party into object
	result.Decode(&party)

	// Return party
	c.JSON(http.StatusOK, party)
}

// Creates a new party
func CreateParty(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the party
	var party models.Party

	// Bind the input from the request body to the party object
	if err := c.ShouldBindJSON(&party); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new obejctId for the party
	party.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(party); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert party
	if _, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").InsertOne(ctx, party); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created party
	c.JSON(http.StatusOK, party)
}

// Changes all fields of a party
func ChangeParty(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the party
	var party models.Party

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the party object
	if err := c.ShouldBindJSON(&party); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	// Set the objectId
	party.Id = objId

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(party); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Replace the existing document with the new one
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").ReplaceOne(ctx, bson.M{"_id": objId}, party)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": result.ModifiedCount,
		"updatedId":     objId,
		"updatedObject": party,
	})
}

// Changes a partys name
func ChangePartyName(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdatePartyNameInput

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Change the name of the party
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"name": input.Name}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the name of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": result.ModifiedCount,
		"updatedId":     objId,
		"updatedName":   input.Name,
	})
}

// Changes a partys abbreviation
func ChangePartyAbbreviation(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdatePartyAbbreviationInput

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Change the abbreviation of the party
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"abbreviation": input.Abbreviation}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the abbreviation of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount":       result.ModifiedCount,
		"updatedId":           objId,
		"updatedAbbreviation": input.Abbreviation,
	})
}

// Changes a partys leader
func ChangePartyLeader(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdatePartyLeaderInput

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Change the leader of the party
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"leader": input.Leader}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the leader of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": result.ModifiedCount,
		"updatedId":     objId,
		"updatedLeader": input.Leader,
	})
}

// Changes a partys logo
func ChangePartyLogo(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdatePartyLogoInput

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Change the logo of the party
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"logo": input.Logo}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the logo of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": result.ModifiedCount,
		"updatedId":     objId,
		"updatedLogo":   input.Logo,
	})
}

// Changes a partys color
func ChangePartyColor(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdatePartyColorInput

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the input object
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Change the color of the party
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"color": input.Color}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the color of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount": result.ModifiedCount,
		"updatedId":     objId,
		"updatedColor":  input.Color,
	})
}

// Deletes a party
func DeleteParty(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Delete the object from the database
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return deleted count and deleted Id
	c.JSON(http.StatusOK, gin.H{
		"deletedCount": result.DeletedCount,
		"deletedId":    objId,
	})
}

// Returns all parties in the collection
func GetParties(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the parties slice
	var parties []models.Party

	// Find all elements in the parties collection
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no party has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no parties",
		})
		return
	}

	// Decode all elements in the database into the party slice
	if err = result.All(ctx, &parties); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, parties)
}

// Returns all specified parties in the collection
func GetPartiesBySlice(c *gin.Context) {
	sliceString := c.Param("slice")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the parties slice
	var parties []models.Party

	// Initialize the input
	input := models.GetSpecificParties{
		Parties: strings.Split(sliceString, "-"),
	}

	// Initialize slice for object ids of the parties
	var partiesIds []primitive.ObjectID

	// loop over parties slice and append object id to object id slice
	for _, party := range input.Parties {
		objId, err := primitive.ObjectIDFromHex(party)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "bad request ids must be in hex",
			})
			return
		}
		partiesIds = append(partiesIds, objId)
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Find all elements in the parties collection
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("parties").Find(ctx, bson.M{"_id": bson.M{"$in": partiesIds}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no party has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no parties",
		})
		return
	}

	// Decode all elements in the database into the party slice
	if err = result.All(ctx, &parties); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return selected parties
	c.JSON(http.StatusOK, parties)
}
