package model

// EntryActivity specifies the activity of a work log entry.
type EntryActivity struct {
	Id          int    // ID of the entry activity
	Description string // Description of the entry activity
}

// NewEntryActivity creates a new EntryActivity model.
func NewEntryActivity() *EntryActivity {
	return &EntryActivity{}
}
