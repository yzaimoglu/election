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

// Get all quarters, important for the updater
func GetQuarters(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the cities
	var quarters []models.Quarter
	quarterName := "quarters"

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(quarterName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &quarters); err == nil {
			c.JSON(http.StatusOK, quarters)
			return
		}
	}*/

	// Sorting by number ascending
	opts := options.Find().SetSort(bson.M{"citynumber": 1})

	// Get quarters
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").Find(ctx, bson.M{}, opts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no quarter has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no quarters",
		})
		return
	}

	// Decode all elements in the database into the quarters slice
	if err = result.All(ctx, &quarters); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	quarterJSON, err := json.Marshal(quarters)
	if err != nil {
		log.Printf("error marshalling " + quarterName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(quarterName, quarterJSON); err != nil {
		log.Printf("error setting " + quarterName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(quarterName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + quarterName + " in the cache: " + err.Error())
	}

	// Return the quarters
	c.JSON(http.StatusOK, quarters)
}

// Get all quarters, important for the updater
func GetQuartersOfDistrict(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	city := c.Param("city")
	district := c.Param("district")
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the cities
	var quarters []models.Quarter
	quarterName := "quarters-" + city + district

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(quarterName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &quarters); err == nil {
			c.JSON(http.StatusOK, quarters)
			return
		}
	}*/

	// Sorting by number ascending
	opts := options.Find().SetSort(bson.M{"citynumber": 1})
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})

	// Get quarters
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").Find(ctx, bson.M{"$and": filter}, opts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no quarter has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no quarters",
		})
		return
	}

	// Decode all elements in the database into the quarters slice
	if err = result.All(ctx, &quarters); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	quarterJSON, err := json.Marshal(quarters)
	if err != nil {
		log.Printf("error marshalling " + quarterName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(quarterName, quarterJSON); err != nil {
		log.Printf("error setting " + quarterName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(quarterName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + quarterName + " in the cache: " + err.Error())
	}

	// Return the quarters
	c.JSON(http.StatusOK, quarters)
}

// Create a quarter
func CreateQuarter(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the quarter
	var quarter models.Quarter

	// Bind the input from the request body to the quarter object
	if err := c.ShouldBindJSON(&quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new objectID for the quarter
	quarter.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert quarter
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").InsertOne(ctx, quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created quarter
	c.JSON(http.StatusOK, quarter)
}

// Get a quarter by its id
func GetQuarterById(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the quarter
	var quarter models.Quarter

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet("quarterwithid-" + id)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &quarter); err == nil {
			c.JSON(http.StatusOK, quarter)
			return
		}
	}

	// ObjectID from id
	var objId primitive.ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "bad request id must be in hex",
		})
		return
	}

	// Find quarter
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").FindOne(ctx, bson.M{"_id": objId})

	// Check if there is a quarter with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no quarter with this id found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	quarterJSON, err := json.Marshal(quarter)
	if err != nil {
		log.Printf("error marshalling quarterwithid-" + id + " to a json object: " + err.Error())
	}
	if err := models.RedisSet("quarterwithid-"+id, quarterJSON); err != nil {
		log.Printf("error setting quarterwithid-" + id + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL("quarterwithid-"+id, 60*5); err != nil {
		log.Printf("error setting ttl for the quarterwithid-" + id + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, quarter)
}

// Get a quarter by its name
func GetQuarterByName(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	name := c.Param("quarter")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the quarter
	var quarter models.Quarter
	quarterName := "quarter-" + city + "-" + district + "-" + name

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(quarterName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &quarter); err == nil {
			c.JSON(http.StatusOK, quarter)
			return
		}
	}

	// Initialize $and input
	var filter []bson.M

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"name": name})

	// Get quarter
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a quarter with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no quarter with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	quarterJSON, err := json.Marshal(quarter)
	if err != nil {
		log.Printf("error marshalling " + quarterName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(quarterName, quarterJSON); err != nil {
		log.Printf("error setting " + quarterName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(quarterName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + quarterName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, quarter)
}

// Change a quarter
func ChangeQuarter(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	name := c.Param("quarter")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the quarter
	var quarter models.Quarter
	var oldQuarter models.Quarter

	// Bind the input from the request body to the quarter object
	if err := c.ShouldBindJSON(&quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Initialize $and input
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"name": name})

	// Update the quarter
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a quarter with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no quarter with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&oldQuarter); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	quarter.Id = oldQuarter.Id

	// Replace object
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").ReplaceOne(ctx, bson.M{"$and": filter}, quarter); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the recently updated quarter
	c.JSON(http.StatusOK, quarter)
}

// Delete a quarter
func DeleteQuarter(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	name := c.Param("quarter")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize $and input
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"name": name})

	// Delete quarter
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("quarters").DeleteOne(ctx, bson.M{"$and": filter})

	// Check if there is a quarter with the filter
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the deleted count
	c.JSON(http.StatusOK, gin.H{
		"deletedCount": result.DeletedCount,
	})
}
