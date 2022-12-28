package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/yzaimoglu/election/parliament/models"
	"github.com/yzaimoglu/election/parliament/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get all constituencies, important for the updater
func GetConstituencies(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the constituencies
	var constituencies []models.District
	constituencyName := "constituencies"

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(constituencyName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &constituencies); err == nil {
			c.JSON(http.StatusOK, constituencies)
			return
		}
	}*/

	// Sorting by number ascending
	opts := options.Find().SetSort(bson.M{"citynumber": 1})

	// Get constituencies
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").Find(ctx, bson.M{}, opts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no constituency has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no constituencies",
		})
		return
	}

	// Decode all elements in the database into the constituencies slice
	if err = result.All(ctx, &constituencies); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	constituencyJSON, err := json.Marshal(constituencies)
	if err != nil {
		log.Printf("error marshalling " + constituencyName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(constituencyName, constituencyJSON); err != nil {
		log.Printf("error setting " + constituencyName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(constituencyName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + constituencyName + " in the cache: " + err.Error())
	}

	// Return the constituencies
	c.JSON(http.StatusOK, constituencies)
}

// Get all constituencies by city, important for the updater
func GetConstituenciesByCity(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	city := c.Param("city")
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the constituencies
	var constituencies []models.District
	constituencyName := "constituencies" + "-" + city

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(constituencyName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &constituencies); err == nil {
			c.JSON(http.StatusOK, constituencies)
			return
		}
	}*/

	// Sorting by number ascending
	opts := options.Find().SetSort(bson.M{"citynumber": 1})

	// Get constituencies
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").Find(ctx, bson.M{"city": city}, opts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no constituency has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no constituencies",
		})
		return
	}

	// Decode all elements in the database into the constituencies slice
	if err = result.All(ctx, &constituencies); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	constituencyJSON, err := json.Marshal(constituencies)
	if err != nil {
		log.Printf("error marshalling " + constituencyName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(constituencyName, constituencyJSON); err != nil {
		log.Printf("error setting " + constituencyName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(constituencyName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + constituencyName + " in the cache: " + err.Error())
	}

	// Return the constituencies
	c.JSON(http.StatusOK, constituencies)
}

// Create a constituency
func CreateConstituency(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the Constituency
	var constituency models.Constituency

	// Bind the input from the request body to the constituency object
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new obejctId for the constituency
	constituency.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert constituency
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").InsertOne(ctx, constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created constituency
	c.JSON(http.StatusOK, constituency)
}

// Get a constituency by its id/name
func GetConstituency(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the constituency
	var constituency models.Constituency

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet("constituency-" + id)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &constituency); err == nil {
			c.JSON(http.StatusOK, constituency)
			return
		}
	}

	// Initialize $or input
	var filter []bson.M
	var objId primitive.ObjectID

	objId, _ = primitive.ObjectIDFromHex(id)
	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})

	// Insert constituency
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").FindOne(ctx, bson.M{"$or": filter})

	// Check if there is a constituency with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no constituency with this id/name found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	constituencyJSON, err := json.Marshal(constituency)
	if err != nil {
		log.Printf("error marshalling constituency-" + id + " to a json object: " + err.Error())
	}
	if err := models.RedisSet("constituency-"+id, constituencyJSON); err != nil {
		log.Printf("error setting constituency-" + id + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL("constituency-"+id, 60*5); err != nil {
		log.Printf("error setting ttl for the constituency-" + id + " in the cache: " + err.Error())
	}

	// Return the recently created constituency
	c.JSON(http.StatusOK, constituency)
}

// Change a constituency
func ChangeConstituency(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the constituency
	var constituency models.Constituency
	var oldConstituency models.Constituency

	// Bind the input from the request body to the city object
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Initialize $or input
	var filter []bson.M
	var objId primitive.ObjectID
	objId, _ = primitive.ObjectIDFromHex(id)

	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})

	// Find old constituency
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").FindOne(ctx, bson.M{"$or": filter})

	// Check if there is a constituency with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no constituency with this id/name found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&oldConstituency); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	constituency.Id = oldConstituency.Id

	// Replace object
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").ReplaceOne(ctx, bson.M{"$or": filter}, constituency); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the recently updated constituency
	c.JSON(http.StatusOK, constituency)
}

// Delete a constituency
func DeleteConstituency(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize $or input
	var filter []bson.M
	var objId primitive.ObjectID

	objId, _ = primitive.ObjectIDFromHex(id)

	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})

	// Delete the constituency
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").DeleteOne(ctx, bson.M{"$or": filter})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return deleted count
	c.JSON(http.StatusOK, gin.H{
		"deletedCount": result.DeletedCount,
	})
}
