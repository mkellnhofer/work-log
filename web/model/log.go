package model

// LogEntries stores data for the log entries view.
type LogEntries struct {
	Summary *LogEntriesSummary
	ListEntries
}

// LogEntriesSummary stores data for the summary in the log entries view.
type LogEntriesSummary struct {
	OvertimeHours         string
	RemainingVacationDays string
}
