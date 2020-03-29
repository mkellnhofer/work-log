package model

// ListDay stores view data for a work day.
type ListDay struct {
	Date          string
	Weekday       string
	ListEntries   []*ListEntry
	WorkDuration  string
	BreakDuration string
}

// NewListDay creates a new ListDay view model.
func NewListDay() *ListDay {
	return &ListDay{}
}
