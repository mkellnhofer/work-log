package model

// ListEntries stores data for the list entries view.
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
