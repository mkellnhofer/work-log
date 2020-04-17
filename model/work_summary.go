package model

import "time"

// WorkSummary stores information about the work of a user.
type WorkSummary struct {
	UserId        int             // ID of the user
	StartTime     time.Time       // Start time
	EndTime       time.Time       // End time
	WorkDurations []*WorkDuration // Work durations (for each entry type)
}

// NewWorkSummary creates a new WorkSummary model.
func NewWorkSummary() *WorkSummary {
	return &WorkSummary{}
}
