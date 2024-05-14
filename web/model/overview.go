package model

// OverviewEntries stores data for the overview entries view.
type OverviewEntries struct {
	CurrMonthName string
	CurrMonth     string
	PrevMonth     string
	NextMonth     string
	Summary       *OverviewEntriesSummary
	Weeks         []*OverviewWeek
	EntriesDays   []*OverviewEntriesDay
}

// OverviewEntriesSummary stores data for the summary in the overview entries view.
type OverviewEntriesSummary struct {
	MonthTargetHours  string
	MonthActualHours  string
	MonthBalanceHours string

	TypePercentages     map[int]int
	RemainingPercentage int
	TypeHours           map[int]string
	RemainingHours      string
}

// OverviewWeek stores view data for a week.
type OverviewWeek struct {
	WeekDays []*OverviewWeekDay
}

// OverviewDay stores view data for a week day.
type OverviewWeekDay struct {
	Date         string
	IsWeekendDay bool
	IsType       map[int]bool
	StartTime    string
	EndTime      string
	Hours        string
	BreakHours   string
}

// OverviewEntriesDay stores view data for a entries day.
type OverviewEntriesDay struct {
	Date         string
	Weekday      string
	IsWeekendDay bool
	Entries      []*OverviewEntry
	Hours        string
}

// OverviewEntry stores view data for a entry.
type OverviewEntry struct {
	IsMissing   bool
	Id          int
	TypeId      int
	Type        string
	StartTime   string
	EndTime     string
	Duration    string
	Activity    string
	Description string
}
