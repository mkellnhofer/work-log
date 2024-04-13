package model

// ListEntries stores data for the list entries view.
type ListEntries struct {
	HasPrevPage bool
	HasNextPage bool
	PrevPageNum int
	NextPageNum int
	Days        []*ListEntriesDay
}

// ListEntriesDay stores view data for a day.
type ListEntriesDay struct {
	Date                         string
	Weekday                      string
	Entries                      []*ListEntry
	WorkDuration                 string
	BreakDuration                string
	WasTargetWorkDurationReached bool
}

// ListEntry stores view data for a entry.
type ListEntry struct {
	IsMissing     bool
	IsOverlapping bool
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	Duration      string
	EntryActivity string
	Description   string
}
