package model

// ListEntry stores view data for a entry.
type ListEntry struct {
	IsMissing     bool
	IsOverlapping bool
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	Duration      string
	EntryActivity string
	Description   string
}

// NewListEntry creates a new ListEntry view model.
func NewListEntry() *ListEntry {
	return &ListEntry{}
}
