package model

// CreateEntry stores data for the create entry view.
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
