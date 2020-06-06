package model

const (
	EntryTypeIdWork     int = 1
	EntryTypeIdTravel   int = 2
	EntryTypeIdVacation int = 3
	EntryTypeIdHoliday  int = 4
	EntryTypeIdIllness  int = 5
)

// EntryType specifies the type of a work log entry.
type EntryType struct {
	Id          int    // ID of the entry type
	Description string // Description of the entry type
}

// NewEntryType creates a new EntryType model.
func NewEntryType(id int, description string) *EntryType {
	return &EntryType{id, description}
}
