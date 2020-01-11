package model

import (
	"time"

	"kellnhofer.com/work-log/constant"
)

// Session stores information about session.
type Session struct {
	Id          string    // ID of the session
	UserId      int       // ID of the user
	ExpireAt    time.Time // Expire time of the session
	PreviousUrl string    // Previous requested URL
}

// NewSession creates a new Session model.
func NewSession(id string) *Session {
	expAt := now().Add(constant.SessionValidity)
	return &Session{id, 0, expAt, ""}
}

// IsExpired returns true if session is expired, otherwise false.
func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(now())
}

// Renew renews the session's expire time.
func (s *Session) Renew() {
	s.ExpireAt = now().Add(constant.SessionValidity)
}

func now() time.Time {
	return time.Now().Truncate(time.Second)
}
