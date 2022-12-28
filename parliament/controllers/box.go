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
)

// Create a ballot box
func CreateBox(c *gin.Context) {
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the box
	var box models.Box

	// Bind the input from the request body to the box object
	if err := c.ShouldBindJSON(&box); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Create a new obejctId for the box
	box.Id = primitive.NewObjectID()

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(box); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Insert box
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").InsertOne(ctx, box); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// Return the recently created box
	c.JSON(http.StatusOK, box)
}

// Get all the boxes by city
func GetBoxesByCity(c *gin.Context) {
	city := c.Param("city")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the boxes
	var boxes []models.Box
	boxName := "boxes-" + city

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(boxName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &boxes); err == nil {
			c.JSON(http.StatusOK, boxes)
			return
		}
	}*/

	// Get box
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").Find(ctx, bson.M{"city": city})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no box has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no boxes",
		})
		return
	}

	// Decode all elements in the database into the box slice
	if err = result.All(ctx, &boxes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(boxes)
	if err != nil {
		log.Printf("error marshalling " + boxName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(boxName, boxJSON); err != nil {
		log.Printf("error setting " + boxName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(boxName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + boxName + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, boxes)
}

// Get all the boxes by constituency
func GetBoxesByConstituency(c *gin.Context) {
	city := c.Param("city")
	constituency := c.Param("constituency")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the boxes
	var boxes []models.Box
	boxName := "boxes-" + city + "-" + constituency

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(boxName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &boxes); err == nil {
			c.JSON(http.StatusOK, boxes)
			return
		}
	}*/

	// Initialize $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"constituency": constituency})

	// Get box
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").Find(ctx, bson.M{"$and": filter})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no box has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no boxes",
		})
		return
	}

	// Decode all elements in the database into the box slice
	if err = result.All(ctx, &boxes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(boxes)
	if err != nil {
		log.Printf("error marshalling " + boxName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(boxName, boxJSON); err != nil {
		log.Printf("error setting " + boxName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(boxName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + boxName + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, boxes)
}

// Get all the boxes by district
func GetBoxesByDistrict(c *gin.Context) {
	city := c.Param("city")
	district := c.Param("district")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the boxes
	var boxes []models.Box
	boxName := "boxes-" + city + "-" + district

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(boxName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &boxes); err == nil {
			c.JSON(http.StatusOK, boxes)
			return
		}
	}*/

	// Initialize $and input
	var filter []bson.M

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})

	// Get box
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").Find(ctx, bson.M{"$and": filter})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no box has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no boxes",
		})
		return
	}

	// Decode all elements in the database into the box slice
	if err = result.All(ctx, &boxes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(boxes)
	if err != nil {
		log.Printf("error marshalling " + boxName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(boxName, boxJSON); err != nil {
		log.Printf("error setting " + boxName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(boxName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + boxName + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, boxes)
}

// Get all the boxes by quarter
func GetBoxesByQuarter(c *gin.Context) {
	city := c.Param("city")
	district := c.Param("district")
	quarter := c.Param("quarter")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the boxes
	var boxes []models.Box
	boxName := "boxes-" + city + "-" + district

	// Check if the result has been cached if so return
	/*redisResult, redisErr := models.RedisGet(boxName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &boxes); err == nil {
			c.JSON(http.StatusOK, boxes)
			return
		}
	}*/

	// Initialize $and input
	var filter []bson.M

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"quarter": quarter})

	// Get box
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").Find(ctx, bson.M{"$and": filter})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return that no box has been found
	if !result.TryNext(ctx) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "there are no boxes",
		})
		return
	}

	// Decode all elements in the database into the box slice
	if err = result.All(ctx, &boxes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(boxes)
	if err != nil {
		log.Printf("error marshalling " + boxName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(boxName, boxJSON); err != nil {
		log.Printf("error setting " + boxName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(boxName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + boxName + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, boxes)
}

// Get a box by its id
func GetBoxById(c *gin.Context) {
	id := c.Param("id")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the box
	var box models.Box

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet("boxwithid-" + id)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &box); err == nil {
			c.JSON(http.StatusOK, box)
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

	// Find box
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").FindOne(ctx, bson.M{"_id": objId})

	// Check if there is a box with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no box with this id found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&box); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(box)
	if err != nil {
		log.Printf("error marshalling boxwithid-" + id + " to a json object: " + err.Error())
	}
	if err := models.RedisSet("boxwithid-"+id, boxJSON); err != nil {
		log.Printf("error setting boxwithid-" + id + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL("boxwithid-"+id, 60*5); err != nil {
		log.Printf("error setting ttl for the boxwithid-" + id + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, box)
}

// Get a box by its number
func GetBoxByNumber(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	number := c.Param("number")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the Box
	var box models.Box
	boxName := "box-" + city + "-" + district + "-" + number

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(boxName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &box); err == nil {
			c.JSON(http.StatusOK, box)
			return
		}
	}

	// Initialize $and input
	var filter []bson.M
	var numberInt int64

	// number to numberInt
	numberInt, _ = strconv.ParseInt(number, 10, 64)

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"number": numberInt})

	// Get box
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a box with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no box with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&box); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Set the result to the cache
	boxJSON, err := json.Marshal(box)
	if err != nil {
		log.Printf("error marshalling " + boxName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(boxName, boxJSON); err != nil {
		log.Printf("error setting " + boxName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(boxName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + boxName + " in the cache: " + err.Error())
	}

	// Return the box
	c.JSON(http.StatusOK, box)
}

// Change a ballot box
func ChangeBox(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	number := c.Param("number")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the box
	var box models.Box
	var oldBox models.Box

	// Bind the input from the request body to the box object
	if err := c.ShouldBindJSON(&box); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Validate the input
	validator := validator.New()
	if err := validator.Struct(box); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	// Initialize $and input
	var filter []bson.M
	var numberInt int64

	numberInt, _ = strconv.ParseInt(number, 10, 64)

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"number": numberInt})

	// Update the box
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a box with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no box with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&oldBox); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	box.Id = oldBox.Id

	// Replace object
	if _, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").ReplaceOne(ctx, bson.M{"$and": filter}, box); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Return the recently updated box
	c.JSON(http.StatusOK, box)
}

// Delete a ballot box
func DeleteBox(c *gin.Context) {
	city := c.Param("id")
	district := c.Param("district")
	number := c.Param("number")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize $and input
	var filter []bson.M
	var numberInt int64

	// number to numberInt
	numberInt, _ = strconv.ParseInt(number, 10, 64)

	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"number": numberInt})

	// Delete box
	result, err := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("boxes").DeleteOne(ctx, bson.M{"$and": filter})

	// Check if there is a box with the filter
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
