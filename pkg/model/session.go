package model

import (
	"time"

	"kellnhofer.com/work-log/pkg/constant"
)

const (
	SessionIdLength = 32
)

// Session stores information about session.
type Session struct {
	Id          string    // ID of the session
	UserId      int       // ID of the user
	ExpireAt    time.Time // Expire time of the session
	PreviousUrl string    // Previous requested URL
}

// NewSession creates a new Session model.
func NewSession() *Session {
	id := generateRandomString(SessionIdLength)
	expAt := now().Add(constant.SessionValidity)
	return &Session{id, AnonymousUserId, expAt, ""}
}

// IsExpired returns true if session is expired, otherwise false.
func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(now())
}

// Renew renews the session's expire time.
func (s *Session) Renew() {
	s.ExpireAt = now().Add(constant.SessionValidity)
}

// GetShortId returns a truncated session ID.
func (s *Session) GetShortId() string {
	return createTruncatedString(s.Id, 8)
}

func IsValidSessionId(sessId string) bool {
	return len(sessId) == SessionIdLength
}
