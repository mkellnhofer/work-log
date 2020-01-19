package model

// CreateEntry stores view data for creating a work entry.
type CreateEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewCreateEntry creates a new CreateEntry view model.
func NewCreateEntry() *CreateEntry {
	return &CreateEntry{}
}
