package model

type SearchQuery interface {
	isSearchQuery()
}

type baseSearchQuery struct {}

func (*baseSearchQuery) isSearchQuery() {}

// BasicSearchQuery stores data for the basic search form view.
type BasicSearchQuery struct {
	baseSearchQuery
	Text string
}

// AdvancedSearchQuery stores data for the advanced search form view.
type AdvancedSearchQuery struct {
	baseSearchQuery
	Input           *AdvancedSearchQueryInput
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// AdvancedSearchQueryInput stores view data of a advanced search form fields.
type AdvancedSearchQueryInput struct {
	ByType         bool
	TypeId         int
	ByDate         bool
	StartDate      string
	StartDateValue string
	EndDate        string
	EndDateValue   string
	ByActivity     bool
	ActivityId     int
	ByDescription  bool
	Description    string
	ByLabels       bool
	Labels         []string
}
