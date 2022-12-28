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

// Loop through the constituencies and do the magic
func LoopConstituency() {
	// Get all the constituencies from the rest api
	constituencies := GetConstituencies()

	// Loop through all the constituencies
	for _, constituency := range constituencies {
		fmt.Println("-----------")
		fmt.Println(constituency.City + "-" + constituency.Name)

		// Initialize candidates of constituency
		oldCandidatesOfConstituency := constituency.Candidates
		var candidatesOfConstituency []models.MVCandidateInBox

		// Loop over the candidates and set the votes to 0 in order to be able to calculate the votes
		for _, candidate := range oldCandidatesOfConstituency {
			candidate.Votes = 0
			candidatesOfConstituency = append(candidatesOfConstituency, candidate)
		}

		// Get all the districts of the district
		districts := GetDistrictsOfConstituency(constituency)
		totalEligibleVoters := int64(0)
		totalValidVotes := int64(0)
		totalInvalidVotes := int64(0)
		totalActualVoters := int64(0)

		// Loop over all the quarters
		for _, district := range districts {
			fmt.Println("District: ", district.Name)
			candidates := district.Candidates
			totalEligibleVoters += district.EligibleVoters
			totalValidVotes += district.ValidVotes
			totalInvalidVotes += district.InvalidVotes
			totalActualVoters += district.ActualVoters

			// Loop over the candidates to be able to calculate
			for candidateIndex, candidate := range candidates {
				// Check if the candidate is corrupt
				if candidatesOfConstituency[candidateIndex].LastName == candidate.LastName {
					// Calculate the new votes
					candidatesOfConstituency[candidateIndex].Votes = candidatesOfConstituency[candidateIndex].Votes + candidate.Votes
					fmt.Println("Candidate: "+candidate.LastName+" (", candidatesOfConstituency[candidateIndex].Votes, ")")
				}
			}
		}
		// Set the new votes with a PUT request to the rest api
		status, statusCode := SetVotesOfConstituency(constituency, candidatesOfConstituency, totalEligibleVoters, totalValidVotes, totalInvalidVotes, totalActualVoters)
		fmt.Println(statusCode, status)
		fmt.Println("-----------")
		fmt.Println("")
		fmt.Println("")
	}
}

// Set the votes of the constituency
func SetVotesOfConstituency(constituency models.MVConstituency, candidates []models.MVCandidateInBox,
	eligibleVoters int64, validVotes int64, invalidVotes int64, actualVoters int64) (string, int) {
	// Initialize the HTTP client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Set the new candidates
	constituency.Candidates = candidates
	constituency.EligibleVoters = eligibleVoters
	constituency.ValidVotes = validVotes
	constituency.InvalidVotes = invalidVotes
	constituency.ActualVoters = actualVoters

	// Marshal the constituency object into JSON
	jsonData, err := json.Marshal(constituency)
	if err != nil {
		log.Fatal(err)
	}

	// Create the PUT request for setting the new candidates
	req, err := http.NewRequest(http.MethodPut, "http://localhost:84/v1/constituency/"+constituency.Id.Hex()+"/", bytes.NewBuffer(jsonData))
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

// Get the districts of constituency
func GetDistrictsOfConstituency(constituency models.MVConstituency) []models.MVDistrict {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the district slice to be returned
	var districts []models.MVDistrict

	// Create the request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/districts/"+constituency.City+"/"+constituency.Name+"/", nil)
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
		fmt.Println("no districts for this constituency found")
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

	// Unmarshal the reponse to the districts object and return
	jsonErr := json.Unmarshal(body, &districts)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return districts
}

// Get the districts
func GetConstituencies() []models.MVConstituency {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the constituencies slice
	var constituencies []models.MVConstituency

	// Create the GET request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/constituencies", nil)
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
		fmt.Println("no constituencies found")
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

	// Unmarshal the json into the constituencies slice
	jsonErr := json.Unmarshal(body, &constituencies)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return constituencies
}
