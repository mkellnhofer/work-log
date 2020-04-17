package model

// CopyEntry stores data for the copy entry view.
type CopyEntry struct {
	PreviousUrl     string
	ErrorMessage    string
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// NewCopyEntry creates a new CopyEntry view model.
func NewCopyEntry() *CopyEntry {
	return &CopyEntry{}
}
