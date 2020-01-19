package model

// EntryType stores view data for a work entry type.
type EntryType struct {
	Id          int
	Description string
}

// NewEntryType creates a new EntryType view model.
func NewEntryType(id int, description string) *EntryType {
	return &EntryType{id, description}
}
