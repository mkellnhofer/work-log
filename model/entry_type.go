package model

// EntryType specifies the type of a work log entry.
type EntryType struct {
	Id          int    // ID of the entry type
	Description string // Description of the entry type
}

// NewEntryType creates a new EntryType model.
func NewEntryType() *EntryType {
	return &EntryType{}
}
