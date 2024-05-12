package model

// LogSummary stores data for the summary in the log view.
type LogSummary struct {
	MonthActualHours string
	MonthTargetHours string

	CurrentLoggedPercent    int
	CurrentRemainingPercent int
	CurrentOvertimePercent  int
	CurrentUndertimePercent int
	CurrentLoggedHours      string
	CurrentRemainingHours   string
	CurrentOvertimeHours    string
	CurrentUndertimeHours   string

	TotalOvertimeHours         string
	TotalRemainingVacationDays string
}
