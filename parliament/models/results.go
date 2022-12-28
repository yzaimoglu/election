package models

// Model for the results
type Result struct {
	Location   string              `json:"location" bson:"location"`
	Candidates []CandidateInResult `json:"candidates" bson:"candidates"`
}

// Model for the candidate in a result
type CandidateInResult struct {
	FirstName  string  `json:"firstname" bson:"firstname"`
	LastName   string  `json:"lastname" bson:"lastname"`
	Percentage float32 `json:"percentage" bson:"percentage"`
}
