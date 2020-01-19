package model

// CopyEntry stores view data for copying a work entry.
type CopyEntry struct {
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewCopyEntry creates a new CopyEntry view model.
func NewCopyEntry() *CopyEntry {
	return &CopyEntry{}
}
