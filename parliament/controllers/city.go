package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/yzaimoglu/election/parliament/models"
	"github.com/yzaimoglu/election/parliament/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get all cities, important for the updater
func GetCities(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the cities
	var cities []models.City
	cityName := "cities"

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(cityName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &cities); err == nil {
			c.JSON(http.StatusOK, cities)
			return
		}
	}*/

	// Sorting by number ascending
	opts := options.Find().SetSort(bson.M{"number": 1})

	// Get cities
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").Find(ctx, bson.M{}, opts)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no city has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no cities",
		})
		return
	}

	// Decode all elements in the database into the cities slice
	if err = result.All(ctx, &cities); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	cityJSON, err := json.Marshal(cities)
	if err != nil {
		log.Printf("error marshalling " + cityName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(cityName, cityJSON); err != nil {
		log.Printf("error setting " + cityName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(cityName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + cityName + " in the cache: " + err.Error())
	}

	// Return the cities
	c.JSON(http.StatusOK, cities)
}

// Create a city
func CreateCity(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the City
	var city models.City

	// Bind the input from the request body to the city object
	if err := c.ShouldBindJSON(&city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new obejctId for the city
	city.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert city
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").InsertOne(ctx, city); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created city
	c.JSON(http.StatusOK, city)
}

// Get a city by its id/name/number
func GetCity(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the City
	var city models.City

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet("city-" + id)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &city); err == nil {
			c.JSON(http.StatusOK, city)
			return
		}
	}

	// Initialize $or input
	var filter []bson.M
	var numberInt int64
	var objId primitive.ObjectID

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		numberInt, _ = strconv.ParseInt(id, 10, 64)
	}

	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})
	filter = append(filter, bson.M{"number": numberInt})

	// Insert city
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").FindOne(ctx, bson.M{"$or": filter})

	// Check if there is a city with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no city with this id/name/number found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&city); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	cityJSON, err := json.Marshal(city)
	if err != nil {
		log.Printf("error marshalling city-" + id + " to a json object: " + err.Error())
	}
	if err := models.RedisSet("city-"+id, cityJSON); err != nil {
		log.Printf("error setting city-" + id + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL("city-"+id, 60*5); err != nil {
		log.Printf("error setting ttl for the city-" + id + " in the cache: " + err.Error())
	}

	// Return the recently created city
	c.JSON(http.StatusOK, city)
}

// Change a city
func ChangeCity(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the City
	var city models.City
	var oldCity models.City

	// Bind the input from the request body to the city object
	if err := c.ShouldBindJSON(&city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(city); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Initialize $or input
	var filter []bson.M
	var objId primitive.ObjectID
	var numberInt int64

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		numberInt, _ = strconv.ParseInt(id, 10, 64)
	}

	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})
	filter = append(filter, bson.M{"number": numberInt})

	// Update the city
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").FindOne(ctx, bson.M{"$or": filter})

	// Check if there is a city with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no city with this id/name/number found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&oldCity); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	city.Id = oldCity.Id

	// Replace object
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").ReplaceOne(ctx, bson.M{"$or": filter}, city); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the recently updated city
	c.JSON(http.StatusOK, city)
}

// Delete a city
func DeleteCity(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize $or input
	var filter []bson.M
	var objId primitive.ObjectID
	var numberInt int64

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		numberInt, _ = strconv.ParseInt(id, 10, 64)
	}

	filter = append(filter, bson.M{"name": id})
	filter = append(filter, bson.M{"_id": objId})
	filter = append(filter, bson.M{"number": numberInt})

	// Delete the city
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").DeleteOne(ctx, bson.M{"$or": filter})
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
