package model

// EntryActivity stores view data for a entry activity.
type EntryActivity struct {
	Id          int
	Description string
}

// NewEntryActivity creates a new EntryActivity view model.
func NewEntryActivity(id int, description string) *EntryActivity {
	return &EntryActivity{id, description}
}
