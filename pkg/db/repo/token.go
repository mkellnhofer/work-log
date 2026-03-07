package repo

import (
	"context"
	"database/sql"
	"fmt"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

// TokenRepo retrieves and stores token related entities.
type TokenRepo struct {
	repo
}

// NewTokenRepo creates a new token repository.
func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{repo{db}}
}

// GetTokensByUserId retrieves all tokens for a user.
func (r *TokenRepo) GetTokensByUserId(ctx context.Context, userId int) ([]*model.Token, error) {
	q := "SELECT id, user_id, name, hashed_token, truncated_token "+
		"FROM token WHERE user_id = ?"

	sh := newTokenScanHelper()
	tokens, qErr := sh.scanRows(r.query(ctx, q, userId))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf(
			"Could not query tokens for user %d from database.", userId), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	return tokens, nil
}

// GetTokenById retrieves a token by its ID.
func (r *TokenRepo) GetTokenByIdAndUserId(ctx context.Context, id int, userId int) (*model.Token,
	error) {
	q := "SELECT id, user_id, name, hashed_token, truncated_token "+
		"FROM token WHERE id = ? AND user_id = ?"

	sh := newTokenScanHelper()
	token, found, qErr := sh.scanRow(r.queryRow(ctx, q, id, userId))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf(
			"Could not read token %d from database.", id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return token, nil
}

// GetTokenByHashedValue retrieves a token by its hashed token value.
func (r *TokenRepo) GetTokenByHashedValue(ctx context.Context, value string) (*model.Token, error) {
	q := "SELECT id, user_id, name, hashed_token, truncated_token "+
		"FROM token WHERE hashed_token = ?"

	sh := newTokenScanHelper()
	token, found, qErr := sh.scanRow(r.queryRow(ctx, q, value))
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not read token from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return token, nil
}

// CreateToken creates a new token.
func (r *TokenRepo) CreateToken(ctx context.Context, token *model.Token) error {
	q := "INSERT INTO token (user_id, name, hashed_token, truncated_token) VALUES (?, ?, ?, ?)"

	id, cErr := r.insert(ctx, q, token.UserId, token.Name, token.HashedToken, token.TruncatedToken)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create token in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}
	token.Id = id
	return nil
}

// DeleteTokenById deletes a token by its ID.
func (r *TokenRepo) DeleteTokenById(ctx context.Context, id int) error {
	q := "DELETE FROM token WHERE id = ?"

	dErr := r.exec(ctx, q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf(
			"Could not delete token %d from database.", id), dErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// --- Scan helper functions ---

func newTokenScanHelper() *scanHelper[*model.Token] {
	return newScanHelper(10, scanTokenFunc)
}

func scanTokenFunc(s scanner) (*model.Token, error) {
	var t model.Token
	err := s.Scan(&t.Id, &t.UserId, &t.Name, &t.HashedToken, &t.TruncatedToken)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
