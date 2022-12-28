package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Model for the party object
type Party struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name         string             `json:"name" bsin:"name,omitempty" validate:"required"`
	Abbreviation string             `json:"abbreviation" bson:"abbreviation,omitempty" validate:"required"`
	Leader       string             `json:"leader" bson:"leader,omitempty" validate:"required"`
	Logo         string             `json:"logo" bson:"logo,omitempty" validate:"required"`
	Color        Color              `json:"color" bson:"color,omitempty" validate:"required"`
}

// Model for the get specific parties
type GetSpecificParties struct {
	Parties []string `json:"parties" validate:"required"`
}

// Model for the update party name input
type UpdatePartyNameInput struct {
	Name string `json:"name" validate:"required"`
}

// Model for the update party abbreviation input
type UpdatePartyAbbreviationInput struct {
	Abbreviation string `json:"abbreviation" validate:"required"`
}

// Model for the update party leader input
type UpdatePartyLeaderInput struct {
	Leader string `json:"leader" validate:"required"`
}

// Model for the update party logo input
type UpdatePartyLogoInput struct {
	Logo string `json:"logo" validate:"required"`
}

// Model for the update party color input
type UpdatePartyColorInput struct {
	Color string `json:"color" validate:"required"`
}
