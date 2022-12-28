package models

// Model for the rank object
type Rank struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Permissions int64  `json:"permissions"`
}
