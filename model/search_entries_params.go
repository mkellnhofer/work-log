package model

import "time"

// SearchEntriesParams stores parameters for a search of work log entries.
type SearchEntriesParams struct {
	ByType        bool      // Flag to search by entry type
	TypeId        int       // ID of the entry type
	ByTime        bool      // Flag to search by start and end time
	StartTime     time.Time // Start time
	EndTime       time.Time // End time
	ByActivity    bool      // Flag to search by entry activity
	ActivityId    int       // ID of the entry activity
	ByDescription bool      // Flag to search by description
	Description   string    // Description
}

// NewSearchEntriesParams create a new SearchEntriesParams model.
func NewSearchEntriesParams() *SearchEntriesParams {
	return &SearchEntriesParams{}
}
