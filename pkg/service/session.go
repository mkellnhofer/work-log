package service

import (
	"context"

	"kellnhofer.com/work-log/pkg/db/repo"
	"kellnhofer.com/work-log/pkg/db/tx"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util"
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
func (s *SessionService) GetSession(ctx context.Context, rawId string) (*model.Session, error) {
	// Hash raw ID
	id := util.CreateHashedString(rawId)

	// Get session by hashed ID
	sess, err := s.sRepo.GetSessionById(ctx, id)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	
	// Set raw ID
	sess.RawId = rawId

	return sess, nil
}

// SaveSession creates/updates a session.
func (s *SessionService) SaveSession(ctx context.Context, session *model.Session) error {
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
func (s *SessionService) DeleteSession(ctx context.Context, rawId string) error {
	// Hash raw ID
	id := util.CreateHashedString(rawId)
	
	// Delete session by hashed ID
	return s.sRepo.DeleteSessionById(ctx, id)
}

// DeleteExpiredSessions deletes expired sessions.
func (s *SessionService) DeleteExpiredSessions(ctx context.Context) error {
	return s.sRepo.DeleteExpiredSessions(ctx)
}
