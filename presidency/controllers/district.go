package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/yzaimoglu/election/presidency/models"
	"github.com/yzaimoglu/election/presidency/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create a district
func CreateDistrict(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the district
	var district models.District

	// Bind the input from the request body to the district object
	if err := c.ShouldBindJSON(&district); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new objectId for the district
	district.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(district); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert district
	if _, err := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("districts").InsertOne(ctx, district); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created district
	c.JSON(http.StatusOK, district)
}

// Get a district by its id
func GetDistrictById(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the district
	var district models.District

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet("districtwithid-" + id)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &district); err == nil {
			c.JSON(http.StatusOK, district)
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

	// Find district
	result := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("districts").FindOne(ctx, bson.M{"_id": objId})

	// Check if there is a district with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no district with this id found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&district); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	districtJSON, err := json.Marshal(district)
	if err != nil {
		log.Printf("error marshalling districtwithid-" + id + " to a json object: " + err.Error())
	}
	if err := models.RedisSet("districtwithid-"+id, districtJSON); err != nil {
		log.Printf("error setting districtwithid-" + id + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL("districtwithid-"+id, 60*5); err != nil {
		log.Printf("error setting ttl for the districtwithid-" + id + " in the cache: " + err.Error())
	}

	// Return the district
	c.JSON(http.StatusOK, district)
}

// Get a district by its name
func GetDistrictByName(c *gin.Context) {
	city := c.Param("id")
	districtParam := c.Param("district")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the district
	var district models.District
	districtName := "district-" + city + "-" + districtParam

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(districtName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &district); err == nil {
			c.JSON(http.StatusOK, district)
			return
		}
	}

	// Initialize $and input
	var filter []bson.M

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"name": districtParam})

	// Get district
	result := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("districts").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a district with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no district with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&district); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	districtJSON, err := json.Marshal(district)
	if err != nil {
		log.Printf("error marshalling " + districtName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(districtName, districtJSON); err != nil {
		log.Printf("error setting " + districtName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(districtName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + districtName + " in the cache: " + err.Error())
	}

	// Return the district
	c.JSON(http.StatusOK, district)
}

// Change a district
func ChangeDistrict(c *gin.Context) {
	city := c.Param("id")
	districtParam := c.Param("district")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the district
	var district models.District
	var oldDistrict models.District

	// Bind the input from the request body to the district object
	if err := c.ShouldBindJSON(&district); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(district); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Initialize $and input
	var filter []bson.M

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": districtParam})

	// Update the district
	result := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("districts").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a district with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no district with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&oldDistrict); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	district.Id = oldDistrict.Id

	// Replace object
	if _, err := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("cities").ReplaceOne(ctx, bson.M{"$and": filter}, district); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the recently updated district
	c.JSON(http.StatusOK, district)
}

// Delete a district
func DeleteDistrict(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize $and input
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})

	// Delete district
	result, err := client.Database(utilities.GetEnv("CB_DB_DATABASE", "cumhurbaskanligi")).Collection("districts").DeleteOne(ctx, bson.M{"$and": filter})

	// Check if there is a district with the filter
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
