package model

import "time"

// Entry stores information about a entry.
type Entry struct {
	Id          int       // ID of the entry
	UserId      int       // ID of the user
	TypeId      int       // ID of the entry type
	StartTime   time.Time // Start time of the entry
	EndTime     time.Time // End time of the entry
	ActivityId  int       // ID of the entry activity
	Project     string    // Related project name of the entry
	Description string    // Description for the entry
	Labels      []string  // Labels for the entry
}

// NewEntry create a new Entry model.
func NewEntry() *Entry {
	return &Entry{}
}
