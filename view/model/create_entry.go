package model

// CreateEntry stores view data for creating a work entry.
type CreateEntry struct {
	ErrorMessage    string
	EntryTypeId     int
	EntryTypes      []*EntryType
	Date            string
	StartTime       string
	EndTime         string
	BreakDuration   string
	EntryActivityId int
	EntryActivities []*EntryActivity
	Description     string
}

// NewCreateEntry creates a new CreateEntry view model.
func NewCreateEntry() *CreateEntry {
	return &CreateEntry{}
}
