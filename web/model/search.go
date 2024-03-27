package model

// Search stores view data of a search.
type Search struct {
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

// NewSearch creates a new Search view model.
func NewSearch() *Search {
	return &Search{}
}
