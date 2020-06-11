package model

import (
	"crypto/rand"
	"encoding/base64"
	"io"
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
func NewSession() *Session {
	id := newSessionId()
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

// --- Helper functions ---

func newSessionId() string {
	b := make([]byte, 24)
	io.ReadFull(rand.Reader, b)
	return base64.URLEncoding.EncodeToString(b)
}

func now() time.Time {
	return time.Now().Truncate(time.Second)
}
