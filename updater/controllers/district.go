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

// Loop through the districts and do the magic
func LoopDistrict() {
	// Get all the districts from the rest api
	districts := GetDistricts()

	// Loop through all the districts
	for _, district := range districts {
		fmt.Println("-----------")
		fmt.Println(district.City + "-" + district.Constituency + "-" + district.Name)

		// Initialize candidates of district
		oldCandidatesOfDistrict := district.Candidates
		var candidatesOfDistrict []models.MVCandidateInBox

		// Loop over the candidates and set the votes to 0 in order to be able to calculate the votes
		for _, candidate := range oldCandidatesOfDistrict {
			candidate.Votes = 0
			candidatesOfDistrict = append(candidatesOfDistrict, candidate)
		}

		// Get all the quarters of the district
		quarters := GetQuartersOfDistrict(district)
		totalEligibleVoters := int64(0)
		totalValidVotes := int64(0)
		totalInvalidVotes := int64(0)
		totalActualVoters := int64(0)

		// Loop over all the quarters
		for _, quarter := range quarters {
			fmt.Println("Quarter: ", quarter.Name)
			candidates := quarter.Candidates
			totalEligibleVoters += quarter.EligibleVoters
			totalValidVotes += quarter.ValidVotes
			totalInvalidVotes += quarter.InvalidVotes
			totalActualVoters += quarter.ActualVoters

			// Loop over the candidates to be able to calculate
			for candidateIndex, candidate := range candidates {
				// Check if the candidate is corrupt
				if candidatesOfDistrict[candidateIndex].LastName == candidate.LastName {
					// Calculate the new votes
					candidatesOfDistrict[candidateIndex].Votes = candidatesOfDistrict[candidateIndex].Votes + candidate.Votes
					fmt.Println("Candidate: "+candidate.LastName+" (", candidatesOfDistrict[candidateIndex].Votes, ")")
				}
			}
		}
		// Set the new votes with a PUT request to the rest api
		status, statusCode := SetVotesOfDistrict(district, candidatesOfDistrict, totalEligibleVoters, totalValidVotes, totalInvalidVotes, totalActualVoters)
		fmt.Println(statusCode, status)
		fmt.Println("-----------")
		fmt.Println("")
		fmt.Println("")
	}
}

// Set the votes of the city
func SetVotesOfDistrict(district models.MVDistrict, candidates []models.MVCandidateInBox,
	eligibleVoters int64, validVotes int64, invalidVotes int64, actualVoters int64) (string, int) {
	// Initialize the HTTP client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Set the new candidates
	district.Candidates = candidates
	district.EligibleVoters = eligibleVoters
	district.ValidVotes = validVotes
	district.InvalidVotes = invalidVotes
	district.ActualVoters = actualVoters

	// Marshal the district object into JSON
	jsonData, err := json.Marshal(district)
	if err != nil {
		log.Fatal(err)
	}

	// Create the PUT request for setting the new candidates
	req, err := http.NewRequest(http.MethodPut, "http://localhost:84/v1/district/"+district.City+"/"+district.Name+"/", bytes.NewBuffer(jsonData))
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

	// Turn the response body into a string and return the result with the status code
	sb := "RESULT: " + string(body)
	return sb, res.StatusCode
}

// Get the quarters of a district
func GetQuartersOfDistrict(district models.MVDistrict) []models.MVQuarter {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the quarter slice to be returned
	var quarters []models.MVQuarter

	// Create the request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/quarters/"+district.City+"/"+district.Name+"/", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the request Headers
	// TODO: AUTHENTICATION
	req.Header.Set("User-Agent", "updater-v1")

	// Execute the request
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	// Return not found if the status code is not 200
	if res.StatusCode != 200 {
		fmt.Println("no quarters for this district found")
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

	// Unmarshal the reponse to the boxes object and return
	jsonErr := json.Unmarshal(body, &quarters)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return quarters
}

// Get the districts
func GetDistricts() []models.MVDistrict {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the districts slice
	var districts []models.MVDistrict

	// Create the GET request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/districts", nil)
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
		fmt.Println("no districts found")
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

	// Unmarshal the json into the districts slice
	jsonErr := json.Unmarshal(body, &districts)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return districts
}
