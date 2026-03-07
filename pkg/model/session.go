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
	Id          string    // Hashed ID of the session (stored in DB)
	RawId       string    // Raw ID of the session (sent as cookie, not stored)
	UserId      int       // ID of the user
	ExpireAt    time.Time // Expire time of the session
	PreviousUrl string    // Previous requested URL
}

// NewSession creates a new Session model.
func NewSession() *Session {
	rawId := generateRandomString(SessionIdLength)
	hashedId := createHashedString(rawId)
	expAt := now().Add(constant.SessionValidity)
	return &Session{
		Id:          hashedId,
		RawId:       rawId,
		UserId:      AnonymousUserId,
		ExpireAt:    expAt,
		PreviousUrl: "",
	}
}

// IsExpired returns true if session is expired, otherwise false.
func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(now())
}

// Renew renews the session's expire time.
func (s *Session) Renew() {
	s.ExpireAt = now().Add(constant.SessionValidity)
}

// GetShortRawId returns a truncated raw session ID.
func (s *Session) GetShortRawId() string {
	return createTruncatedString(s.RawId, 8)
}

func IsValidSessionId(sessId string) bool {
	return len(sessId) == SessionIdLength
}
