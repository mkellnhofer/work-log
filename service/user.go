package service

import (
	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/db/repo"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/model"
)

// UserService contains user related logic.
type UserService struct {
	uRepo *repo.UserRepo
}

// NewUserService create a new user service.
func NewUserService(ur *repo.UserRepo) *UserService {
	return &UserService{ur}
}

// --- User functions ---

// GetUserById gets a user by its ID.
func (s *UserService) GetUserById(id int) (*model.User, *e.Error) {
	return s.uRepo.GetUserById(id)
}

// GetUserByUsername gets a user by its username.
func (s *UserService) GetUserByUsername(username string) (*model.User, *e.Error) {
	return s.uRepo.GetUserByUsername(username)
}

// --- User settings functions ---

// GetSettingShowOverviewDetails gets the setting value for the "show overview details" setting.
func (s *UserService) GetSettingShowOverviewDetails(userId int) (bool, *e.Error) {
	return s.uRepo.GetUserBoolSetting(userId, constant.SettingKeyShowOverviewDetails)
}

// UpdateSettingShowOverviewDetails updates the setting value for the "show overview details" setting.
func (s *UserService) UpdateSettingShowOverviewDetails(userId int, value bool) *e.Error {
	return s.uRepo.UpdateUserBoolSetting(userId, constant.SettingKeyShowOverviewDetails, value)
}

// --- User contract functions ---

// GetUserContractByUserId gets the contract information of a user by its ID.
func (s *UserService) GetUserContractByUserId(userId int) (*model.UserContract, *e.Error) {
	return s.uRepo.GetUserContractByUserId(userId)
}
