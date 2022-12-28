package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yzaimoglu/election/updater/models"
)

// Loop through the cities and do the magic
func LoopCity() {
	// Get all the cities from the rest api
	cities := GetCities()

	// Loop through all the cities
	for _, city := range cities {
		fmt.Println("-----------")
		fmt.Println(city.Name)

		// Initialize candidates of the city
		oldCandidatesOfCity := city.Candidates
		var candidatesOfCity []models.MVCandidateInBox

		// Loop over the candidates and set the votes to 0 in order to be able to calculate the votes
		for _, candidate := range oldCandidatesOfCity {
			candidate.Votes = 0
			candidatesOfCity = append(candidatesOfCity, candidate)
		}

		// Get all the constituencies of the city
		constituencies := GetConstituenciesOfCity(city)
		totalEligibleVoters := int64(0)
		totalValidVotes := int64(0)
		totalInvalidVotes := int64(0)
		totalActualVoters := int64(0)

		// Loop over all the constituencies
		for _, constituency := range constituencies {
			fmt.Println("Name: " + constituency.Name)
			candidates := constituency.Candidates
			totalEligibleVoters += constituency.EligibleVoters
			totalValidVotes += constituency.ValidVotes
			totalInvalidVotes += constituency.InvalidVotes
			totalActualVoters += constituency.ActualVoters

			// Loop over the candidates to be able to calculate
			for candidateIndex, candidate := range candidates {
				// Check if candidate is corrupt
				if candidatesOfCity[candidateIndex].LastName == candidate.LastName {
					// Calculate the new votes
					candidatesOfCity[candidateIndex].Votes = candidatesOfCity[candidateIndex].Votes + candidate.Votes
					fmt.Println("Candidate: "+candidate.LastName+" (", candidatesOfCity[candidateIndex].Votes, ")")
				}
			}
		}
		// Set the new votes with a PUT request to the rest api
		status, statusCode := SetVotesOfCity(city, candidatesOfCity, totalEligibleVoters, totalValidVotes, totalInvalidVotes, totalActualVoters)
		fmt.Println(statusCode, status)
		fmt.Println("-----------")
		fmt.Println("")
		fmt.Println("")
	}
}

// Set the votes of the city
func SetVotesOfCity(city models.MVCity, candidates []models.MVCandidateInBox,
	eligibleVoters int64, validVotes int64, invalidVotes int64, actualVoters int64) (string, int) {
	// Initialize the HTTP client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Set the new candidates
	city.Candidates = candidates
	city.EligibleVoters = eligibleVoters
	city.ValidVotes = validVotes
	city.InvalidVotes = invalidVotes
	city.ActualVoters = actualVoters

	// Marshal the city object into JSON
	jsonData, err := json.Marshal(city)
	if err != nil {
		log.Fatal(err)
	}

	// Create the PUT request for setting the new candidates
	req, err := http.NewRequest(http.MethodPut, "http://localhost:84/v1/city/"+city.Name, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// Set the request headers
	// TODO: AUTHENTICATION
	req.Header.Set("User-Agent", "updater-v1")

	// Execute the recently created request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Get the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Turn the response body into a string and return the results with the status code
	sb := "RESULT: " + string(body)
	return sb, res.StatusCode
}

// Get the boxes of a city
func GetConstituenciesOfCity(city models.MVCity) []models.MVConstituency {
	// Initialize the HTTP client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the constituencies slice to be returned
	var constituencies []models.MVConstituency

	// Create the request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/constituencies/"+city.Name, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the request headers
	// TODO: AUTHENTICATION
	req.Header.Set("User-Agent", "updater-v1")

	// Execute the request
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Return not found if the status code is not 200
	if res.StatusCode != 200 {
		fmt.Println("no constituencies for this city found")
		return nil
	}

	// Close the response body
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Read the body
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Unmarshal the response to the constituencies object and return
	jsonErr := json.Unmarshal(body, &constituencies)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return constituencies
}

// Get an object from the json
func GetCities() []models.MVCity {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the cities slice
	var cities []models.MVCity

	// Creat the GET request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/cities", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the request headers
	// TODO: AUTHENTICATION
	req.Header.Set("User-Agent", "updater-v1")

	// Execute the request
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Return not found if the status code is not 200
	if res.StatusCode != 200 {
		fmt.Println("no cities found")
		return nil
	}

	// Close the response body
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Read the body
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// Unmarshal the json into the cities slice
	jsonErr := json.Unmarshal(body, &cities)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return cities
}
