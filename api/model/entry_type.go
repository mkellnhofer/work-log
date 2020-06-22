package model

// EntryType
//
// Specifies the type of a entry.
//
// swagger:model EntryType
type EntryType struct {
	// The ID of the entry type.
	// example: 1
	Id int `json:"id"`

	// The description of the entry type.
	// min length: 1
	// max length: 50
	// example: Work
	Description string `json:"description"`
}
