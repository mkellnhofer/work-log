package model

// SearchEntries stores view data for searching work entries.
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
