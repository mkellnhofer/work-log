package service

import (
	"context"

	"kellnhofer.com/work-log/db/repo"
	"kellnhofer.com/work-log/db/tx"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/model"
)

// SessionService contains session related logic.
type SessionService struct {
	service
	sRepo *repo.SessionRepo
}

// NewSessionService create a new session service.
func NewSessionService(tm *tx.TransactionManager, sr *repo.SessionRepo) *SessionService {
	return &SessionService{service{tm}, sr}
}

// --- Session functions ---

// GetSession gets a session.
func (s *SessionService) GetSession(ctx context.Context, id string) (*model.Session, *e.Error) {
	return s.sRepo.GetSessionById(ctx, id)
}

// SaveSession creates/updates a session.
func (s *SessionService) SaveSession(ctx context.Context, session *model.Session) *e.Error {
	// Check if session exists
	exists, err := s.sRepo.ExistsSessionById(ctx, session.Id)
	if err != nil {
		return err
	}

	// Exists session?
	if exists {
		// Update session
		return s.sRepo.UpdateSession(ctx, session)
	} else {
		// Create session
		return s.sRepo.CreateSession(ctx, session)
	}
}

// DeleteSession deletes a session.
func (s *SessionService) DeleteSession(ctx context.Context, id string) *e.Error {
	return s.sRepo.DeleteSessionById(ctx, id)
}

// DeleteExpiredSessions deletes expired sessions.
func (s *SessionService) DeleteExpiredSessions(ctx context.Context) *e.Error {
	return s.sRepo.DeleteExpiredSessions(ctx)
}
