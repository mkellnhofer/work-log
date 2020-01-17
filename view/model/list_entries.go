package model

// ListEntries stores view data for listing work entries.
type ListEntries struct {
	HasPrevPage bool
	HasNextPage bool
	PrevPageNum int
	NextPageNum int
	Days        []*Day
}

// NewListEntries creates a new ListEntries view model.
func NewListEntries() *ListEntries {
	return &ListEntries{}
}

// Day stores view data for a work day.
type Day struct {
	Date         string
	Weekday      string
	Entries      []*Entry
	WorkDuration string
}

// NewDay creates a new Day view model.
func NewDay() *Day {
	return &Day{}
}

// Entry stores view data for a work entry.
type Entry struct {
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	BreakDuration string
	WorkDuration  string
	EntryActivity string
	Description   string
}

// NewEntry creates a new Entry view model.
func NewEntry() *Entry {
	return &Entry{}
}
