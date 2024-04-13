package model

// OverviewEntries stores data for the overview entries view.
type OverviewEntries struct {
	PreviousUrl   string
	CurrMonthName string
	CurrMonth     string
	PrevMonth     string
	NextMonth     string
	Summary       *OverviewEntriesSummary
	Days          []*OverviewEntriesDay
}

// OverviewEntriesSummary stores data for the summary in the overview entries view.
type OverviewEntriesSummary struct {
	ActualWorkHours     string
	ActualTravelHours   string
	ActualVacationHours string
	ActualHolidayHours  string
	ActualIllnessHours  string
	TargetHours         string
	ActualHours         string
	BalanceHours        string
}

// OverviewEntriesDay stores view data for a day.
type OverviewEntriesDay struct {
	Date         string
	Weekday      string
	IsWeekendDay bool
	Entries      []*OverviewEntry
	WorkDuration string
}

// OverviewEntry stores view data for a entry.
type OverviewEntry struct {
	Id            int
	EntryType     string
	StartTime     string
	EndTime       string
	Duration      string
	EntryActivity string
	Description   string
}
