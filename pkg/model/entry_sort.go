package model

// EntrySort stores parameters to sort entries.
type EntrySort struct {
	ByTime Sorting // Time sorting direction
}

// NewEntrySort create a new EntrySort model.
func NewEntrySort() *EntrySort {
	return &EntrySort{}
}
