package model

// UpdateEntryActivity
//
// Holds the new information about a entry activity.
//
// swagger:model UpdateEntryActivity
type UpdateEntryActivity struct {
	// The description of the entry activity.
	// min length: 1
	// max length: 50
	// example: Development
	Description string `json:"description"`
}
