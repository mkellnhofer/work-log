package service

import (
	"context"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/db/repo"
	"kellnhofer.com/work-log/db/tx"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/model"
)

// UserService contains user related logic.
type UserService struct {
	service
	uRepo *repo.UserRepo
}

// NewUserService create a new user service.
func NewUserService(tm *tx.TransactionManager, ur *repo.UserRepo) *UserService {
	return &UserService{service{tm}, ur}
}

// --- User functions ---

// GetUserById gets a user by its ID.
func (s *UserService) GetUserById(ctx context.Context, id int) (*model.User, *e.Error) {
	return s.uRepo.GetUserById(ctx, id)
}

// GetUserByUsername gets a user by its username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User,
	*e.Error) {
	return s.uRepo.GetUserByUsername(ctx, username)
}

// --- User settings functions ---

// GetSettingShowOverviewDetails gets the setting value for the "show overview details" setting.
func (s *UserService) GetSettingShowOverviewDetails(ctx context.Context, userId int) (bool,
	*e.Error) {
	return s.uRepo.GetUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails)
}

// UpdateSettingShowOverviewDetails updates the setting value for the "show overview details" setting.
func (s *UserService) UpdateSettingShowOverviewDetails(ctx context.Context, userId int,
	value bool) *e.Error {
	return s.uRepo.UpdateUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails, value)
}

// --- User contract functions ---

// GetUserContractByUserId gets the contract information of a user by its ID.
func (s *UserService) GetUserContractByUserId(ctx context.Context, userId int) (*model.UserContract,
	*e.Error) {
	return s.uRepo.GetUserContractByUserId(ctx, userId)
}
