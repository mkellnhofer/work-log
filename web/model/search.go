package model

// SearchQuery stores data for the search form view.
type SearchQuery struct {
	Input           *SearchQueryInput
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// SearchQueryInput stores view data of a search form fields.
type SearchQueryInput struct {
	ByType         bool
	TypeId         int
	ByDate         bool
	StartDate      string
	StartDateValue string
	EndDate        string
	EndDateValue   string
	ByActivity     bool
	ActivityId     int
	ByLabels       bool
	Labels         []string
	Text           string
}

// SearchDetails stores data for the search details view.
type SearchDetails struct {
	EntryFilterDetails
}
