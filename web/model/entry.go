package model

// Entry stores view data of a entry.
type Entry struct {
	Id          int
	TypeId      int
	Date        string
	StartTime   string
	EndTime     string
	ActivityId  int
	Description string
}

// CreateEntry stores data for the create entry view.
type CreateEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// CopyEntry stores data for the copy entry view.
type CopyEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// EditEntry stores data for the edit entry view.
type EditEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// DeleteEntry stores data for the delete entry view.
type DeleteEntry struct {
	PreviousUrl  string
	ErrorMessage string
	Id           int
}
