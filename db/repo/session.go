package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

type dbSession struct {
	id          string
	userId      sql.NullInt64
	expireAt    string
	previousUrl sql.NullString
}

// SessionRepo retrieves and stores sessions related entities.
type SessionRepo struct {
	repo
}

// NewSessionRepo creates a new session repository.
func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{repo{db}}
}

// --- Session functions ---

// GetSessionById retrieves a session by its ID.
func (r *SessionRepo) GetSessionById(ctx context.Context, id string) (*model.Session, *e.Error) {
	sr, qErr := r.queryRow(ctx, &scanSessionHelper{}, "SELECT id, user_id, expire_at, previous_url "+
		"FROM session WHERE id = ?", id)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read session %s from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.Session), nil
}

// ExistsSessionById checks if a session exists.
func (r *SessionRepo) ExistsSessionById(ctx context.Context, id string) (bool, *e.Error) {
	cnt, cErr := r.count(ctx, "session", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read session %s from "+
			"database.", id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateSession creates a new session.
func (r *SessionRepo) CreateSession(ctx context.Context, session *model.Session) *e.Error {
	sess := toDbSession(session)

	cErr := r.exec(ctx, "INSERT INTO session (id, user_id, expire_at, previous_url) "+
		"VALUES (?, ?, ?, ?)", sess.id, sess.userId, sess.expireAt, sess.previousUrl)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create session in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// UpdateSession updates a session.
func (r *SessionRepo) UpdateSession(ctx context.Context, session *model.Session) *e.Error {
	sess := toDbSession(session)

	uErr := r.exec(ctx, "UPDATE session SET user_id = ?, expire_at = ?, previous_url = ? WHERE id = ?",
		sess.userId, sess.expireAt, sess.previousUrl, sess.id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update session %s in database.",
			session.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteSessionById deletes a session by by its ID.
func (r *SessionRepo) DeleteSessionById(ctx context.Context, id string) *e.Error {
	dErr := r.exec(ctx, "DELETE FROM session WHERE id = ?", id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete session %s from "+
			"database.", id), dErr)
		log.Error(err.StackTrace())
	}

	return nil
}

// DeleteExpiredSessions deletes expired sessions.
func (r *SessionRepo) DeleteExpiredSessions(ctx context.Context) *e.Error {
	now := time.Now()
	n := *formatTimestamp(&now)
	dErr := r.exec(ctx, "DELETE FROM session WHERE expire_at < ?", n)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, "Could not delete expired sessions from database.",
			dErr)
		log.Error(err.StackTrace())
	}

	return nil
}

// --- Helper functions ---

type scanSessionHelper struct {
}

func (h *scanSessionHelper) makeSlice() interface{} {
	return make([]*model.Session, 0, 10)
}

func (h *scanSessionHelper) scan(s scanner) (interface{}, error) {
	var dbS dbSession

	err := s.Scan(&dbS.id, &dbS.userId, &dbS.expireAt, &dbS.previousUrl)
	if err != nil {
		return nil, err
	}

	session := fromDbSession(&dbS)

	return session, nil
}

func (h *scanSessionHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.Session), item.(*model.Session))
}

func toDbSession(in *model.Session) *dbSession {
	var out dbSession
	out.id = in.Id
	if in.UserId != 0 {
		out.userId = sql.NullInt64{Int64: int64(in.UserId), Valid: true}
	} else {
		out.userId = sql.NullInt64{Int64: 0, Valid: false}
	}
	out.expireAt = *formatTimestamp(&in.ExpireAt)
	if in.PreviousUrl != "" {
		out.previousUrl = sql.NullString{String: in.PreviousUrl, Valid: true}
	} else {
		out.previousUrl = sql.NullString{String: "", Valid: false}
	}
	return &out
}

func fromDbSession(in *dbSession) *model.Session {
	var out model.Session
	out.Id = in.id
	if in.userId.Valid {
		out.UserId = int(in.userId.Int64)
	} else {
		out.UserId = 0
	}
	out.ExpireAt = *parseTimestamp(&in.expireAt)
	if in.previousUrl.Valid {
		out.PreviousUrl = in.previousUrl.String
	} else {
		out.PreviousUrl = ""
	}
	return &out
}
