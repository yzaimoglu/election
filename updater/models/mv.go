package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the ballot box object
type MVBox struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Number         int64              `json:"number" bson:"number"`             // 1001
	City           string             `json:"city" bson:"city"`                 // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber"`     // 6
	Constituency   string             `json:"constituency" bson:"constituency"` // Ankara-1
	District       string             `json:"district" bson:"district"`         // Çankaya
	Quarter        string             `json:"quarter" bson:"quarter"`           // Çukurambar
	Candidates     []MVCandidateInBox `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
	SST            string             `json:"sst" bson:"sst"`                       // 24923948264 (Static File Storage Microservice)
	SDC            string             `json:"sdc" bson:"sdc"`                       // 42424234242 (Static File Storage Microservice)
}

// Model for the city
type MVCity struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name" validate:"required"`             // Ankara
	Number         int64              `json:"number" bson:"number" validate:"required,numeric"` // 6
	Candidates     []MVCandidateInBox `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters" validate:"numeric"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters" validate:"numeric"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes" validate:"numeric"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes" validate:"numeric"`     // 161
}

// Model for the constituency
type MVConstituency struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name" validate:"required"`
	City           string             `json:"city" bson:"city" validate:"required"`                     // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber" validate:"required,numeric"` // 6
	Candidates     []MVCandidateInBox `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
}

// Model for the district
type MVDistrict struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`                 // Cankaya
	City           string             `json:"city" bson:"city"`                 // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber"`     // 6
	Constituency   string             `json:"constituency" bson:"constituency"` // ankara-1
	Candidates     []MVCandidateInBox `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
}

// Model for the quarter
type MVQuarter struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`                 // Cevizlidere
	City           string             `json:"city" bson:"city"`                 // Ankara
	CityNumber     int64              `json:"citynumber" bson:"citynumber"`     // 6
	Constituency   string             `json:"constituency" bson:"constituency"` // Ankara-1
	District       string             `json:"district" bson:"district"`         // Cankaya
	Candidates     []MVCandidateInBox `json:"candidates" bson:"candidates"`
	EligibleVoters int64              `json:"eligiblevoters" bson:"eligiblevoters"` // 12621
	ActualVoters   int64              `json:"actualvoters" bson:"actualvoters"`     // 10262
	ValidVotes     int64              `json:"validvotes" bson:"validvotes"`         // 10101
	InvalidVotes   int64              `json:"invalidvotes" bson:"invalidvotes"`     // 161
}

// Model for a Party in a Box
type MVCandidateInBox struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	Votes     int64  `json:"votes" bson:"votes"` // 121
}
