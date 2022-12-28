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

// Loop through the quarters and do the magic
func LoopQuarter() {
	// Get all the quarters from the rest api
	quarters := GetQuarters()

	// Loop through all the quarters
	for _, quarter := range quarters {
		fmt.Println("-----------")
		fmt.Println(quarter.City + "-" + quarter.District + "-" + quarter.Name)

		// Initialize candidates of quarter
		oldCandidatesOfQuarter := quarter.Candidates
		var candidatesOfQuarter []models.MVCandidateInBox

		// Loop over the candidates and set the votes to 0 in order to be able to calculate the votes
		for _, candidate := range oldCandidatesOfQuarter {
			candidate.Votes = 0
			candidatesOfQuarter = append(candidatesOfQuarter, candidate)
		}

		// Get all the boxes of the quarter
		boxes := GetBoxesOfQuarter(quarter)
		totalEligibleVoters := int64(0)
		totalValidVotes := int64(0)
		totalInvalidVotes := int64(0)
		totalActualVoters := int64(0)

		// Loop over all the boxes
		for _, box := range boxes {
			fmt.Println("Box: ", box.Number)
			candidates := box.Candidates
			totalEligibleVoters += box.EligibleVoters
			totalValidVotes += box.ValidVotes
			totalInvalidVotes += box.InvalidVotes
			totalActualVoters += box.ActualVoters

			// Loop over the candidates to be able to calculate
			for candidateIndex, candidate := range candidates {
				// Check if the candidate is corrupt
				if candidatesOfQuarter[candidateIndex].LastName == candidate.LastName {
					// Calculate the new votes
					candidatesOfQuarter[candidateIndex].Votes = candidatesOfQuarter[candidateIndex].Votes + candidate.Votes
					fmt.Println("Candidate: "+candidate.LastName+" (", candidatesOfQuarter[candidateIndex].Votes, ")")
				}
			}
		}
		// Set the new votes with a PUT request to the rest api
		status, statusCode := SetVotesOfQuarter(quarter, candidatesOfQuarter, totalEligibleVoters, totalValidVotes, totalInvalidVotes, totalActualVoters)
		fmt.Println(statusCode, status)
		fmt.Println("-----------")
		fmt.Println("")
		fmt.Println("")
	}
}

// Set the votes of the city
func SetVotesOfQuarter(quarter models.MVQuarter, candidates []models.MVCandidateInBox,
	eligibleVoters int64, validVotes int64, invalidVotes int64, actualVoters int64) (string, int) {
	// Initialize the HTTP client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Set the new candidates
	quarter.Candidates = candidates
	quarter.EligibleVoters = eligibleVoters
	quarter.ValidVotes = validVotes
	quarter.InvalidVotes = invalidVotes
	quarter.ActualVoters = actualVoters

	// Marshal the quarter object into JSON
	jsonData, err := json.Marshal(quarter)
	if err != nil {
		log.Fatal(err)
	}

	// Create the PUT request for setting the new candidates
	req, err := http.NewRequest(http.MethodPut, "http://localhost:84/v1/quarter/"+quarter.City+"/"+quarter.District+"/"+quarter.Name+"/", bytes.NewBuffer(jsonData))
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

// Get the boxes of a quarter
func GetBoxesOfQuarter(quarter models.MVQuarter) []models.MVBox {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the box slice to be returned
	var boxes []models.MVBox

	// Create the request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/boxes/"+quarter.City+"/"+quarter.District+"/"+quarter.Name, nil)
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
		fmt.Println("no boxes for this quarter found")
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
	jsonErr := json.Unmarshal(body, &boxes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return boxes
}

// Get the quarters
func GetQuarters() []models.MVQuarter {
	// Initialize the HTTP Client
	client := http.Client{
		Timeout: time.Second * 15,
	}

	// Initialize the quarters slice
	var quarters []models.MVQuarter

	// Create the GET request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:84/v1/quarters", nil)
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
		fmt.Println("no quarters found")
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

	// Unmarshal the json into the quarters slice
	jsonErr := json.Unmarshal(body, &quarters)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return quarters
}
