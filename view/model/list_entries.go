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

// ListDay stores view data for a work day.
type ListDay struct {
	Date         string
	Weekday      string
	ListEntries  []*ListEntry
	WorkDuration string
}

// NewListDay creates a new ListDay view model.
func NewListDay() *ListDay {
	return &ListDay{}
}

// ListEntry stores view data for a work entry.
type ListEntry struct {
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	BreakDuration string
	WorkDuration  string
	EntryActivity string
	Description   string
}

// NewListEntry creates a new ListEntry view model.
func NewListEntry() *ListEntry {
	return &ListEntry{}
}
