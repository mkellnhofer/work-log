package constant

import (
	"time"
)

const (
	AppVersion string = "1.0.0-alpha"

	DbTimestampFormat string = "2006-01-02 15:04:05"

	SessionValidity time.Duration = 1 * time.Hour

	EntryTypeWork     int = 1
	EntryTypeTravel   int = 2
	EntryTypeVacation int = 3
	EntryTypeHoliday  int = 4
	EntryTypeIllness  int = 5
)
