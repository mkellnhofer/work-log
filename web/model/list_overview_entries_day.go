package model

// ListOverviewEntriesDay stores view data for a day.
type ListOverviewEntriesDay struct {
	Date         string
	Weekday      string
	IsWeekendDay bool
	Entries      []*ListOverviewEntry
	WorkDuration string
}

// NewListOverviewEntriesDay creates a new ListOverviewEntriesDay view model.
func NewListOverviewEntriesDay() *ListOverviewEntriesDay {
	return &ListOverviewEntriesDay{}
}
