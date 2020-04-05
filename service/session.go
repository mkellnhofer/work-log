package service

import (
	"kellnhofer.com/work-log/db/repo"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/model"
)

// SessionService contains session related logic.
type SessionService struct {
	sRepo *repo.SessionRepo
}

// NewSessionService create a new session service.
func NewSessionService(sr *repo.SessionRepo) *SessionService {
	return &SessionService{sr}
}

// --- Session functions ---

// GetSession gets a session.
func (s *SessionService) GetSession(id string) (*model.Session, *e.Error) {
	return s.sRepo.GetSessionById(id)
}

// SaveSession creates/updates a session.
func (s *SessionService) SaveSession(session *model.Session) *e.Error {
	// Check if session exists
	exists, err := s.sRepo.ExistsSessionById(session.Id)
	if err != nil {
		return err
	}

	// Exists session?
	if exists {
		// Update session
		return s.sRepo.UpdateSession(session)
	} else {
		// Create session
		return s.sRepo.CreateSession(session)
	}
}

// DeleteSession deletes a session.
func (s *SessionService) DeleteSession(id string) *e.Error {
	return s.sRepo.DeleteSessionById(id)
}

// DeleteExpiredSessions deletes expired sessions.
func (s *SessionService) DeleteExpiredSessions() *e.Error {
	return s.sRepo.DeleteExpiredSessions()
}
