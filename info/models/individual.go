package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the individual object
type Individual struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName   string             `json:"firstname" bson:"firstname,omitempty" validate:"required"`
	LastName    string             `json:"lastname" bson:"lastname,omitempty" validate:"required"`
	BirthDate   string             `json:"birthdate" bson:"birthdate,omitempty" validate:"required"`
	Image       string             `json:"image" bson:"image" validate:"required"`
	Affiliation string             `json:"affiliation" bson:"affiliation"`
	Color       Color              `json:"color" bson:"color" validate:"required"`
}

// Model for the get specific individuals
type GetSpecificIndividuals struct {
	Individuals []string `json:"individuals" validate:"required"`
}

// Model for the update individual ƒirstname input
type UpdateIndividualFirstNameInput struct {
	FirstName string `json:"firstname" validate:"required"`
}

// Model for the update individual lastname input
type UpdateIndividualLastNameInput struct {
	LastName string `json:"lastname" validate:"required"`
}

// Model for the update party leader input
type UpdateIndividualBirthDateInput struct {
	BirthDate string `json:"birthdate" validate:"required"`
}
