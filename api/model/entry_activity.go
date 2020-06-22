package model

// EntryActivity
//
// Specifies the activity of a entry.
//
// swagger:model EntryActivity
type EntryActivity struct {
	// The ID of the entry activity.
	// example: 1
	Id int `json:"id"`

	// The description of the entry activity.
	// min length: 1
	// max length: 50
	// example: Development
	Description string `json:"description"`
}
