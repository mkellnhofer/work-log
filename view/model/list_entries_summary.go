package model

// ListEntriesSummary stores data for the summary in the list entries view.
type ListEntriesSummary struct {
	OvertimeHours         string
	RemainingVacationDays string
}

// NewListEntriesSummary creates a new ListEntriesSummary view model.
func NewListEntriesSummary() *ListEntriesSummary {
	return &ListEntriesSummary{}
}
