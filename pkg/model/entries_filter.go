package model

import "time"

// EntriesFilter stores parameters to filter entries.
type EntriesFilter struct {
	ByUser        bool      // Flag to filter by user
	UserId        int       // ID of the user
	ByType        bool      // Flag to filter by entry type
	TypeId        int       // ID of the entry type
	ByTime        bool      // Flag to filter by time
	StartTime     time.Time // Start time
	EndTime       time.Time // End time
	ByActivity    bool      // Flag to filter by entry activity
	ActivityId    int       // ID of the entry activity
	ByDescription bool      // Flag to filter by description
	Description   string    // Description
}

// NewEntriesFilter create a new EntriesFilter model.
func NewEntriesFilter() *EntriesFilter {
	return &EntriesFilter{}
}
