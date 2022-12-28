package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the city
type City struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name" validate:"required"`             // Ankara
	Number         int64              `json:"number" bson:"number" validate:"required,numeric"` // 6
	Parties        []PartyInBox       `json:"parties" bson:"parties"`
	Individuals    []IndividualInBox  `json:"individuals" bson:"individuals"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters" validate:"numeric"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters" validate:"numeric"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes" validate:"numeric"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes" validate:"numeric"`     // 161
}
