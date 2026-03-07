package service

import (
	"context"
	"fmt"

	"kellnhofer.com/work-log/pkg/db/repo"
	"kellnhofer.com/work-log/pkg/db/tx"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util"
)

// TokenService contains token related logic.
type TokenService struct {
	service
	tRepo *repo.TokenRepo
}

// NewTokenService create a new token service.
func NewTokenService(tm *tx.TransactionManager, tr *repo.TokenRepo) *TokenService {
	return &TokenService{service{tm}, tr}
}

// --- Token functions ---

// GetTokenByValue gets a token by its token value.
func (s *TokenService) GetTokenByValue(ctx context.Context, value string) (*model.Token, error) {
	hashedToken := util.CreateHashedString(value)
	return s.tRepo.GetTokenByHashedValue(ctx, hashedToken)
}

// --- Current user token functions ---

// CreateCurrentUserToken creates a new token for the current user.
func (s *TokenService) CreateCurrentUserToken(ctx context.Context, name string) (*model.Token,
	error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserAccount); err != nil {
		return nil, err
	}

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Create new token
	token := model.NewToken(userId, name)

	// Store token
	if err := s.tRepo.CreateToken(ctx, token); err != nil {
		return nil, err
	}
	
	return token, nil
}

// GetCurrentUserTokens gets all tokens of the current user.
func (s *TokenService) GetCurrentUserTokens(ctx context.Context) ([]*model.Token, error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserAccount); err != nil {
		return nil, err
	}

	// Get tokens
	return s.tRepo.GetTokensByUserId(ctx, getCurrentUserId(ctx))
}

// GetCurrentUserTokenById gets a token by ID for the current user.
func (s *TokenService) GetCurrentUserTokenById(ctx context.Context, id int) (*model.Token,
	error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserAccount); err != nil {
		return nil, err
	}

	// Get token
	return s.getCurrentUserTokenById(ctx, id)
}

// DeleteCurrentUserTokenById deletes a token by ID for the current user.
func (s *TokenService) DeleteCurrentUserTokenById(ctx context.Context, id int) error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserAccount); err != nil {
		return err
	}

	// Get token
	token, err := s.getCurrentUserTokenById(ctx, id)
	if err != nil {
		return err
	}

	// Remove token
	return s.tRepo.DeleteTokenById(ctx, token.Id)
}

func (s *TokenService) getCurrentUserTokenById(ctx context.Context, id int) (*model.Token, error) {
	token, err := s.tRepo.GetTokenByIdAndUserId(ctx, id, getCurrentUserId(ctx))
	if err != nil {
		return nil, err
	}
	if token == nil {
		err := e.NewError(e.LogicTokenNotFound, fmt.Sprintf("Could not find token %d.", id))
		log.Debug(err.StackTrace())
		return nil, err
	}
	return token, nil
}
