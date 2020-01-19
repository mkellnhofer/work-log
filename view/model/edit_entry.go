package model

// EditEntry stores view data for editing a work entry.
type EditEntry struct {
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewEditEntry creates a new EditEntry view model.
func NewEditEntry() *EditEntry {
	return &EditEntry{}
}
