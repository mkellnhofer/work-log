package model

import "time"

// WorkDuration stores the work and break durations for a entry type.
type WorkDuration struct {
	TypeId        int           // ID of the entry type
	WorkDuration  time.Duration // Work duration
	BreakDuration time.Duration // Break duration
}

// NewWorkDuration creates a new WorkDuration model.
func NewWorkDuration() *WorkDuration {
	return &WorkDuration{}
}
