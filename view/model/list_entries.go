package model

// ListEntries stores data for the list entries view.
type ListEntries struct {
	Summary     *ListEntriesSummary
	HasPrevPage bool
	HasNextPage bool
	PrevPageNum int
	NextPageNum int
	Days        []*ListEntriesDay
}

// NewListEntries creates a new ListEntries view model.
func NewListEntries() *ListEntries {
	return &ListEntries{}
}
