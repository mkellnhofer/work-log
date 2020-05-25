package model

// ListOverviewEntriesSummary stores data for the summary in the list overview entries view.
type ListOverviewEntriesSummary struct {
	WorkDescription     string
	TravelDescription   string
	VacationDescription string
	HolidayDescription  string
	IllnessDescription  string
	ActualWorkHours     string
	ActualTravelHours   string
	ActualVacationHours string
	ActualHolidayHours  string
	ActualIllnessHours  string
	TargetHours         string
	ActualHours         string
	BalanceHours        string
}

// NewListOverviewEntriesSummary creates a new ListOverviewEntriesSummary view model.
func NewListOverviewEntriesSummary() *ListOverviewEntriesSummary {
	return &ListOverviewEntriesSummary{}
}
