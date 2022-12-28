package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yzaimoglu/election/parliament/models"
	"github.com/yzaimoglu/election/parliament/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get the results by city
func GetResultsByCity(c *gin.Context) {
	city := c.Param("city")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the result
	var resultObj models.Result
	var cityObj models.City
	resultName := "results-" + city
	resultObj.Location = resultName

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(resultName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &resultObj); err == nil {
			c.JSON(http.StatusOK, resultObj)
			return
		}
	}

	// Get results
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("cities").FindOne(ctx, bson.M{"name": city})

	// Check if there is a city with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no city with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&cityObj); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Loop over the candidates and add calculate the percentages
	candidates := cityObj.Candidates
	total := cityObj.ValidVotes
	var candidatesOfResult []models.CandidateInResult
	for _, candidate := range candidates {
		percentage := float32(candidate.Votes) / float32(total) * 100
		newCandidate := models.CandidateInResult{
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Percentage: percentage,
		}
		candidatesOfResult = append(candidatesOfResult, newCandidate)
	}
	resultObj.Candidates = candidatesOfResult

	// Set the result to the cache
	resultJSON, err := json.Marshal(resultObj)
	if err != nil {
		log.Printf("error marshalling " + resultName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(resultName, resultJSON); err != nil {
		log.Printf("error setting " + resultName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(resultName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + resultName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, resultObj)
}

// Get the results by city
func GetResultsByConstituency(c *gin.Context) {
	city := c.Param("city")
	constituency := c.Param("constituency")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the result
	var resultObj models.Result
	var constituencyObj models.Constituency
	resultName := "results-" + city + "-" + constituency
	resultObj.Location = resultName

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(resultName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &resultObj); err == nil {
			c.JSON(http.StatusOK, resultObj)
			return
		}
	}

	// Initialize the $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"name": constituency})

	// Get results
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("constituencies").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a constituency with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no constituency with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&constituencyObj); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Loop over the candidates and add calculate the percentages
	candidates := constituencyObj.Candidates
	total := constituencyObj.ValidVotes
	var candidatesOfResult []models.CandidateInResult
	for _, candidate := range candidates {
		percentage := float32(candidate.Votes) / float32(total) * 100
		newCandidate := models.CandidateInResult{
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Percentage: percentage,
		}
		candidatesOfResult = append(candidatesOfResult, newCandidate)
	}
	resultObj.Candidates = candidatesOfResult

	// Set the result to the cache
	resultJSON, err := json.Marshal(resultObj)
	if err != nil {
		log.Printf("error marshalling " + resultName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(resultName, resultJSON); err != nil {
		log.Printf("error setting " + resultName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(resultName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + resultName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, resultObj)
}

// Get the results by city
func GetResultsByDistrict(c *gin.Context) {
	city := c.Param("city")
	constituency := c.Param("constituency")
	district := c.Param("district")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the result
	var resultObj models.Result
	var districtObj models.District
	resultName := "results-" + city + "-" + constituency + "-" + district
	resultObj.Location = resultName

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(resultName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &resultObj); err == nil {
			c.JSON(http.StatusOK, resultObj)
			return
		}
	}

	// Initialize the $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"constituency": constituency})
	filter = append(filter, bson.M{"name": district})

	// Get results
	result := client.Database(utilities.GetEnv("MV_DB_DATABASE", "milletvekili")).Collection("districts").FindOne(ctx, bson.M{"$and": filter})

	// Check if there is a district with the filter
	if result.Err() == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no district with these details found",
		})
		return
	}

	// Decode result to object
	if err := result.Decode(&districtObj); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Loop over the candidates and add calculate the percentages
	candidates := districtObj.Candidates
	total := districtObj.ValidVotes
	var candidatesOfResult []models.CandidateInResult
	for _, candidate := range candidates {
		percentage := float32(candidate.Votes) / float32(total) * 100
		newCandidate := models.CandidateInResult{
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Percentage: percentage,
		}
		candidatesOfResult = append(candidatesOfResult, newCandidate)
	}
	resultObj.Candidates = candidatesOfResult

	// Set the result to the cache
	resultJSON, err := json.Marshal(resultObj)
	if err != nil {
		log.Printf("error marshalling " + resultName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(resultName, resultJSON); err != nil {
		log.Printf("error setting " + resultName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(resultName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + resultName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, resultObj)
}

// Get the results by city
func GetResultsByQuarter(c *gin.Context) {
	city := c.Param("city")
	constituency := c.Param("constituency")
	district := c.Param("district")
	quarter := c.Param("quarter")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the result
	var resultObj models.Result
	var quarterObj models.Quarter
	resultName := "results-" + city + "-" + constituency + "-" + district + "-" + quarter
	resultObj.Location = resultName

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(resultName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &resultObj); err == nil {
			c.JSON(http.StatusOK, resultObj)
			return
		}
	}

	// Initialize the $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"constituency": constituency})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"name": quarter})

	// Get results
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
	if err := result.Decode(&quarterObj); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Loop over the candidates and add calculate the percentages
	candidates := quarterObj.Candidates
	total := quarterObj.ValidVotes
	var candidatesOfResult []models.CandidateInResult
	for _, candidate := range candidates {
		percentage := float32(candidate.Votes) / float32(total) * 100
		newCandidate := models.CandidateInResult{
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Percentage: percentage,
		}
		candidatesOfResult = append(candidatesOfResult, newCandidate)
	}
	resultObj.Candidates = candidatesOfResult

	// Set the result to the cache
	resultJSON, err := json.Marshal(resultObj)
	if err != nil {
		log.Printf("error marshalling " + resultName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(resultName, resultJSON); err != nil {
		log.Printf("error setting " + resultName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(resultName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + resultName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, resultObj)
}

// Get the results by city
func GetResultsByBox(c *gin.Context) {
	city := c.Param("city")
	constituency := c.Param("constituency")
	district := c.Param("district")
	quarter := c.Param("quarter")
	box := c.Param("box")
	client, ctx, cancel := models.GetMongoInstance()
	defer cancel()
	defer client.Disconnect(ctx)

	// Initialize the result
	var resultObj models.Result
	var boxObj models.Box
	resultName := "results-" + city + "-" + constituency + "-" + district + "-" + quarter + "-" + box
	resultObj.Location = resultName

	// Check if the result has been cached if so return
	redisResult, redisErr := models.RedisGet(resultName)
	if redisErr == nil {
		if err := json.Unmarshal(redisResult, &resultObj); err == nil {
			c.JSON(http.StatusOK, resultObj)
			return
		}
	}

	// Parse the box into a boxnumber
	boxNumber, err := strconv.ParseInt(box, 0, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "the box number must be a number",
		})
		return
	}

	// Initialize the $and filter
	var filter []bson.M
	filter = append(filter, bson.M{"city": city})
	filter = append(filter, bson.M{"constituency": constituency})
	filter = append(filter, bson.M{"district": district})
	filter = append(filter, bson.M{"quarter": quarter})
	filter = append(filter, bson.M{"number": boxNumber})

	// Get results
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
	if err := result.Decode(&boxObj); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "internal server error",
		})
		return
	}

	// Loop over the candidates and add calculate the percentages
	candidates := boxObj.Candidates
	total := boxObj.ValidVotes
	var candidatesOfResult []models.CandidateInResult
	for _, candidate := range candidates {
		percentage := float32(candidate.Votes) / float32(total) * 100
		newCandidate := models.CandidateInResult{
			FirstName:  candidate.FirstName,
			LastName:   candidate.LastName,
			Percentage: percentage,
		}
		candidatesOfResult = append(candidatesOfResult, newCandidate)
	}
	resultObj.Candidates = candidatesOfResult

	// Set the result to the cache
	resultJSON, err := json.Marshal(resultObj)
	if err != nil {
		log.Printf("error marshalling " + resultName + " to a json object: " + err.Error())
	}
	if err := models.RedisSet(resultName, resultJSON); err != nil {
		log.Printf("error setting " + resultName + " to the cache: " + err.Error())
	}
	if err := models.RedisTTL(resultName, 60*5); err != nil {
		log.Printf("error setting ttl for the " + resultName + " in the cache: " + err.Error())
	}

	// Return the quarter
	c.JSON(http.StatusOK, resultObj)
}
