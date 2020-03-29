package model

// ListEntry stores view data for a work entry.
type ListEntry struct {
	IsMissing     bool
	IsOverlapping bool
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	BreakDuration string
	WorkDuration  string
	EntryActivity string
	Description   string
}

// NewListEntry creates a new ListEntry view model.
func NewListEntry() *ListEntry {
	return &ListEntry{}
}
