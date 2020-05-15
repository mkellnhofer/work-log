package model

// ListEntriesDay stores view data for a day.
type ListEntriesDay struct {
	Date                         string
	Weekday                      string
	Entries                      []*ListEntry
	WorkDuration                 string
	BreakDuration                string
	WasTargetWorkDurationReached string
}

// NewListEntriesDay creates a new ListEntriesDay view model.
func NewListEntriesDay() *ListEntriesDay {
	return &ListEntriesDay{}
}
