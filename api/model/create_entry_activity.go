package model

// CreateEntryActivity
//
// Holds information about a new entry activity.
//
// swagger:model CreateEntryActivity
type CreateEntryActivity struct {
	// The description of the entry activity.
	// min length: 1
	// max length: 50
	// example: Development
	Description string `json:"description"`
}
