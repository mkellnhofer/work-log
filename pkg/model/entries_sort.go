package model

// EntriesSort stores parameters to sort entries.
type EntriesSort struct {
	ByTime Sorting // Time sorting direction
}

// NewEntriesSort create a new EntriesSort model.
func NewEntriesSort() *EntriesSort {
	return &EntriesSort{}
}
