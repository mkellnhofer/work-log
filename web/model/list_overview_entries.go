package model

// ListOverviewEntries stores data for the list overview entries view.
type ListOverviewEntries struct {
	PreviousUrl   string
	CurrMonthName string
	CurrMonth     string
	PrevMonth     string
	NextMonth     string
	Summary       *ListOverviewEntriesSummary
	ShowDetails   bool
	Days          []*ListOverviewEntriesDay
}

// NewListOverviewEntries creates a new ListOverviewEntries view model.
func NewListOverviewEntries() *ListOverviewEntries {
	return &ListOverviewEntries{}
}
