package model

// SearchEntries stores data for the search entries view.
type SearchEntries struct {
	PreviousUrl     string
	ErrorMessage    string
	ByType          bool
	TypeId          int
	ByDate          bool
	StartDate       string
	EndDate         string
	ByActivity      bool
	ActivityId      int
	ByDescription   bool
	Description     string
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewSearchEntries creates a new SearchEntries view model.
func NewSearchEntries() *SearchEntries {
	return &SearchEntries{}
}
