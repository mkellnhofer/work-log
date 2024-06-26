package constant

import (
	"time"
)

const (
	AppVersion string = "1.3.0-alpha"

	DbDateFormat      string = "2006-01-02"
	DbTimestampFormat string = "2006-01-02 15:04:05"

	ApiDateFormat      string = "2006-01-02"
	ApiTimestampFormat string = "2006-01-02T15:04:05"

	SessionCookieName string        = "session"
	SessionValidity   time.Duration = 1 * time.Hour

	ContextKeyTransactionHolder contextKey = contextKey("transaction-holder")
	ContextKeySessionHolder     contextKey = contextKey("session-holder")
	ContextKeySecurityContext   contextKey = contextKey("security-context")

	ViewPathDefault string = "/log"

	ApiPath string = "/api/v1"
)

type contextKey string

func (c contextKey) String() string {
	return "work log context key " + string(c)
}
