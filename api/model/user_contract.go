package model

// UserContract contains information about the work contract of a user.
type UserContract struct {
	DailyWorkingDuration float32 `json:"dailyWorkingDuration"`
	AnnualVacationDays   float32 `json:"annualVacationDays"`
	InitOvertimeDuration float32 `json:"initOvertimeDuration"`
	InitVacationDays     float32 `json:"initVacationDays"`
	FirstWorkDay         string  `json:"firstWorkDay"`
}
