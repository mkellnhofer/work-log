package model

// EntryData stores data for the create/copy/edit entry view.
type EntryData struct {
	Entry           *Entry
	EntryTypes      []*EntryType
	EntryActivities []*EntryActivity
}

// Entry stores view data of a entry.
type Entry struct {
	Id          int
	TypeId      int
	Date        string
	StartTime   string
	EndTime     string
	ActivityId  int
	Description string
}

const (
	EntryTypeIdWork     int = 1
	EntryTypeIdTravel   int = 2
	EntryTypeIdVacation int = 3
	EntryTypeIdHoliday  int = 4
	EntryTypeIdIllness  int = 5
)

// EntryTypeColors defines colors for each entry type.
var EntryTypeColors = map[int]string{
	EntryTypeIdWork:     "#0c63e4",
	EntryTypeIdTravel:   "#6aa6ff",
	EntryTypeIdVacation: "#04a17a",
	EntryTypeIdHoliday:  "#12dfab",
	EntryTypeIdIllness:  "#e8571e",
}

// EntryType stores view data for a entry type.
type EntryType struct {
	Id          int
	Description string
}

// EntryActivity stores view data for a entry activity.
type EntryActivity struct {
	Id          int
	Description string
}
