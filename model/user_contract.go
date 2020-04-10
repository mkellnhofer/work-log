package model

import "time"

// UserContract stores information about the work contract of a user.
type UserContract struct {
	DailyWorkingDuration time.Duration // Daily working duration of the user
	AnnualVacationDays   float32       // Annual vacation days of the user
	InitOvertimeDuration time.Duration // Initial overtime duration of the user
	InitVacationDays     float32       // Initial vacation days of the user
	FirstWorkDay         time.Time     // First work day of the user
}

// NewUserContract creates a new UserContract model.
func NewUserContract() *UserContract {
	return &UserContract{}
}
