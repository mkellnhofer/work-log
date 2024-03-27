package model

// SearchEntries stores data for the search entries view.
type SearchEntries struct {
	PreviousUrl     string
	ErrorMessage    string
	Search          *Search
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewSearchEntries creates a new SearchEntries view model.
func NewSearchEntries() *SearchEntries {
	return &SearchEntries{}
}
