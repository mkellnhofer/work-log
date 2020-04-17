package model

// ListDay stores view data for a day.
type ListDay struct {
	Date                         string
	Weekday                      string
	ListEntries                  []*ListEntry
	WorkDuration                 string
	BreakDuration                string
	WasTargetWorkDurationReached string
}

// NewListDay creates a new ListDay view model.
func NewListDay() *ListDay {
	return &ListDay{}
}
