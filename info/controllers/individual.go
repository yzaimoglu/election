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

// Returns a single individual with a specific id
func GetIndividual(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	var individual models.Individual

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Find the individual in the database
	result := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").FindOne(ctx, &bson.M{"_id": objId})

	// Check if there is an individual with that id
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no individual with that id found",
		})
		return
	}

	// Decode individual into object
	result.Decode(&individual)

	// Return individual
	c.JSON(http.StatusOK, individual)
}

// Creates a new individual
func CreateIndividual(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the individual
	var individual models.Individual

	// Bind the input from the request body to the individual object
	if err := c.ShouldBindJSON(&individual); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new obejctId for the individual
	individual.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(individual); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert individual
	if _, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").InsertOne(ctx, individual); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created individual
	c.JSON(http.StatusOK, individual)
}

// Changes all fields of an indiviudal
func ChangeIndividual(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the individual
	var individual models.Individual

	// Object ID from id param
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id should be in hex",
		})
		return
	}

	// Bind the input from the request body to the individual object
	if err := c.ShouldBindJSON(&individual); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	// Set the objectId
	individual.Id = objId

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(individual); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Replace the existing document with the new one
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").ReplaceOne(ctx, bson.M{"_id": objId}, individual)
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
		"updatedObject": individual,
	})
}

// Changes an individuals first name
func ChangeIndividualFirstName(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdateIndividualFirstNameInput

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

	// Change the first name of the individual
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"firstname": input.FirstName}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the first name of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount":    result.ModifiedCount,
		"updatedId":        objId,
		"updatedFirstName": input.FirstName,
	})
}

// Changes an individuals last name
func ChangeIndividualLastName(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdateIndividualLastNameInput

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

	// Change the last name of the individual
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"lastname": input.LastName}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the last name of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount":   result.ModifiedCount,
		"updatedId":       objId,
		"updatedLastName": input.LastName,
	})
}

// Changes an individuals birthdate
func ChangeIndividualBirthdate(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the input object
	var input models.UpdateIndividualBirthDateInput

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

	// Change the birth date of the individual
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"birthdate": input.BirthDate}})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the id and the first name of the updated object
	c.JSON(http.StatusOK, gin.H{
		"modifiedCount":    result.ModifiedCount,
		"updatedId":        objId,
		"updatedBirthDate": input.BirthDate,
	})
}

// Deletes an individual
func DeleteIndividual(c *gin.Context) {
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
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").DeleteOne(ctx, bson.M{"_id": objId})
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

// Returns all individuals in the collection
func GetIndividuals(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	var individuals []models.Individual

	// Find all elements in teh individuals collection
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no individual has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no individuals",
		})
		return
	}

	// Decode all elements in the database into the individual slice
	if err = result.All(ctx, &individuals); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, individuals)
}

// Returns all specified individuals in the collection
func GetIndividualsBySlice(c *gin.Context) {
	sliceString := c.Param("slice")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the individuals slice
	var individuals []models.Individual

	// Decode the base64 string
	fromBase64String := string(utilities.FromBase64(sliceString))

	// Get the slice from the concatenated string
	individualsSlice := strings.Split(fromBase64String, "+")

	// Initialize the first and last names
	var individualFirstNames []string
	var individualLastNames []string

	// Loop over the inidividuals
	for _, individualInSlice := range individualsSlice {
		firstAndLastName := strings.Split(individualInSlice, "-")
		individualFirstNames = append(individualFirstNames, firstAndLastName[0])
		individualLastNames = append(individualLastNames, firstAndLastName[1])
	}

	// Initialize the $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"firstname": bson.M{"$in": individualFirstNames}})
	filter = append(filter, bson.M{"lastname": bson.M{"$in": individualLastNames}})

	// Find all elements in the individuals collection
	result, err := client.Database(utilities.GetEnv("BILGI_DB_DATABASE", "bilgi")).Collection("individuals").Find(ctx, bson.M{"$and": filter})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no individual has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no individuals",
		})
		return
	}

	// Decode all elements in the database into the individuals slice
	if err = result.All(ctx, &individuals); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return selected individuals
	c.JSON(http.StatusOK, individuals)
}
