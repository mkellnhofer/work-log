package model

// UserContract
//
// Contains information about the work contract of a user.
//
// swagger:model UserContract
type UserContract struct {

	// The daily working duration of the user in hours.
	// example: 8.0
	DailyWorkingDuration float32 `json:"dailyWorkingDuration"`

	// The annual vacation days of the user.
	// example: 26.0
	AnnualVacationDays float32 `json:"annualVacationDays"`

	// The initial overtime duration of the user in hours.
	// example: 4.5
	InitOvertimeDuration float32 `json:"initOvertimeDuration"`

	// The initial vacation days of the user.
	// example: 2.5
	InitVacationDays float32 `json:"initVacationDays"`

	// The first work day of the user.
	// example: 2019-01-01
	FirstWorkDay string `json:"firstWorkDay"`
}
