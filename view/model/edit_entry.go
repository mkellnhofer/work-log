package model

// EditEntry stores data for the edit entry view.
type EditEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewEditEntry creates a new EditEntry view model.
func NewEditEntry() *EditEntry {
	return &EditEntry{}
}
