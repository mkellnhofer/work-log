package model

// Search stores data for the search view.
type Search struct {
	ErrorMessage    string
	SearchInput     *SearchInput
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// SearchInput stores view data of a search.
type SearchInput struct {
	ByType        bool
	TypeId        int
	ByDate        bool
	StartDate     string
	EndDate       string
	ByActivity    bool
	ActivityId    int
	ByDescription bool
	Description   string
}

// SearchEntries stores data for the search entries view.
type SearchEntries struct {
	Query string
	ListEntries
}
