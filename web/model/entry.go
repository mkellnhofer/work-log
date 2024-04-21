package model

// EntryData stores data for the create/copy/edit entry view.
type EntryData struct {
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

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

// EntryType stores view data for a entry type.
type EntryType struct {
	Id          int
	Description string
}

// EntryActivity stores view data for a entry activity.
type EntryActivity struct {
	Id          int
	Description string
}
