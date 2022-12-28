package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the constituency
type Constituency struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name" validate:"required"`                     // ankara-1
	ReadableName   string             `json:"readablename" bson:"readablename"`                         // Ankara-01
	City           string             `json:"city" bson:"city" validate:"required"`                     // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber" validate:"required,numeric"` // 6
	Candidates     []CandidateInBox   `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
}
