package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the ballot box object
type Box struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Number         int64              `json:"number" bson:"number"`             // 1001
	City           string             `json:"city" bson:"city"`                 // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber"`     // 6
	Constituency   string             `json:"constituency" bson:"constituency"` // Ankara-1
	District       string             `json:"district" bson:"district"`         // Çankaya
	Quarter        string             `json:"quarter" bson:"quarter"`           // Çukurambar
	Candidates     []CandidateInBox   `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
	SST            string             `json:"sst" bson:"sst"`                       // 24923948264 (Static File Storage Microservice)
	SDC            string             `json:"sdc" bson:"sdc"`                       // 42424234242 (Static File Storage Microservice)
}

// Model for a Party in a Box
type CandidateInBox struct {
	FirstName string `json:"firstname" bson:"firstname"` // Recep Tayyip
	LastName  string `json:"lastname" bson:"lastname"`   // Erdogan
	Votes     int64  `json:"votes" bson:"votes"`         // 121
}
