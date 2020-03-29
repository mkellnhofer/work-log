package constant

import (
	"time"
)

const (
	AppVersion string = "1.0.0-alpha"

	DbDateFormat      string = "2006-01-02"
	DbTimestampFormat string = "2006-01-02 15:04:05"

	SessionCookieName string        = "session"
	SessionValidity   time.Duration = 1 * time.Hour

	ContextKeySession contextKey = contextKey("session")

	PathDefault string = "/list"

	EntryTypeWork     int = 1
	EntryTypeTravel   int = 2
	EntryTypeVacation int = 3
	EntryTypeHoliday  int = 4
	EntryTypeIllness  int = 5
)

type contextKey string

func (c contextKey) String() string {
	return "work log context key " + string(c)
}
