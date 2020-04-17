package model

// Entry stores view data of a entry.
type Entry struct {
	Id            int
	TypeId        int
	Date          string
	StartTime     string
	EndTime       string
	BreakDuration string
	ActivityId    int
	Description   string
}

// NewEntry creates a new Entry view model.
func NewEntry() *Entry {
	return &Entry{}
}
