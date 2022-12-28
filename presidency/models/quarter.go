package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the quarter
type Quarter struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`                 // Cevizlidere
	City           string             `json:"city" bson:"city"`                 // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber"`     // 6
	Constituency   string             `json:"constituency" bson:"constituency"` // Ankara-1
	District       string             `json:"district" bson:"district"`         // Cankaya
	Parties        []PartyInBox       `json:"parties" bson:"parties"`
	Individuals    []IndividualInBox  `json:"individuals" bson:"individuals"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
}
