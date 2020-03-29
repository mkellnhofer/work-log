package model

// ListSearchEntries stores view data for listing searched work entries.
type ListSearchEntries struct {
	PreviousUrl string
	Query       string
	HasPrevPage bool
	HasNextPage bool
	PrevPageNum int
	NextPageNum int
	ListDays    []*ListDay
}

// NewListSearchEntries creates a new ListSearchEntries view model.
func NewListSearchEntries() *ListSearchEntries {
	return &ListSearchEntries{}
}
