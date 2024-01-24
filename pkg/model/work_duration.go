package model

import "time"

// WorkDuration stores the work durations for a entry type.
type WorkDuration struct {
	TypeId       int           // ID of the entry type
	WorkDuration time.Duration // Work duration
}

// NewWorkDuration creates a new WorkDuration model.
func NewWorkDuration() *WorkDuration {
	return &WorkDuration{}
}
