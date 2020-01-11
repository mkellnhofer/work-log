package model

import "time"

// Entry stores information about a work log entry.
type Entry struct {
	Id            int           // ID of the entry
	UserId        int           // ID of the user
	TypeId        int           // ID of the entry type
	StartTime     time.Time     // Start time of the entry
	EndTime       time.Time     // End time of the entry
	BreakDuration time.Duration // Break duration of the entry
	ActivityId    int           // ID of the entry activity
	Description   string        // Description for the entry
}

// NewEntry create a new Entry model.
func NewEntry() *Entry {
	return &Entry{}
}
