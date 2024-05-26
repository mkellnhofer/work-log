package model

// SearchQuery stores data for the search form view.
type SearchQuery struct {
	IsAdvanced      bool
	Input           *SearchQueryInput
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// SearchQueryInput stores view data of a search form fields.
type SearchQueryInput struct {
	ByType     bool
	TypeId     int
	ByDate     bool
	StartDate  string
	EndDate    string
	ByActivity bool
	ActivityId int
	Text       string
}

// SearchEntries stores data for the search entries view.
type SearchEntries struct {
	Query string
	ListEntries
}
