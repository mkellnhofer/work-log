package model

// ListEntries stores view data for listing work entries.
type ListEntries struct {
	HasPrevPage bool
	HasNextPage bool
	PrevPageNum int
	NextPageNum int
	ListDays    []*ListDay
}

// NewListEntries creates a new ListEntries view model.
func NewListEntries() *ListEntries {
	return &ListEntries{}
}
