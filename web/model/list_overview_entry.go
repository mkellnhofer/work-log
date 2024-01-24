package model

// ListOverviewEntry stores view data for a entry.
type ListOverviewEntry struct {
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	Duration      string
	EntryActivity string
	Description   string
}

// NewListOverviewEntry creates a new ListOverviewEntry view model.
func NewListOverviewEntry() *ListOverviewEntry {
	return &ListOverviewEntry{}
}
