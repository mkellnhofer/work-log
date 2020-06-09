package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/db/repo"
	"kellnhofer.com/work-log/db/tx"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
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

// GetUsers gets all users.
func (s *UserService) GetUsers(ctx context.Context) ([]*model.User, *e.Error) {
	return s.uRepo.GetUsers(ctx)
}

// GetUserById gets a user by its ID.
func (s *UserService) GetUserById(ctx context.Context, id int) (*model.User, *e.Error) {
	return s.uRepo.GetUserById(ctx, id)
}

// GetUserByUsername gets a user by its username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User,
	*e.Error) {
	return s.uRepo.GetUserByUsername(ctx, username)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, user *model.User) *e.Error {
	// Check if username is already taken
	if err := s.checkIfUsernameIsAlreadyTaken(ctx, 0, user.Username); err != nil {
		return err
	}

	// Hash password
	hashUserPassword(user)

	// Create user
	return s.uRepo.CreateUser(ctx, user)
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) *e.Error {
	// Check if user exists
	if err := s.checkIfUserExists(ctx, user.Id); err != nil {
		return err
	}

	// Check if username is already taken
	if err := s.checkIfUsernameIsAlreadyTaken(ctx, user.Id, user.Username); err != nil {
		return err
	}

	// Hash password
	hashUserPassword(user)

	// Update user
	return s.uRepo.UpdateUser(ctx, user)
}

func hashUserPassword(user *model.User) {
	pPassword := user.Password
	hBytes, hErr := bcrypt.GenerateFromPassword([]byte(pPassword), bcrypt.DefaultCost)
	if hErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not hash user password.", hErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	hPassword := string(hBytes)
	user.Password = hPassword
}

// DeleteUserById deletes a user by its ID.
func (s *UserService) DeleteUserById(ctx context.Context, id int) *e.Error {
	// Check if user exists
	if err := s.checkIfUserExists(ctx, id); err != nil {
		return err
	}

	// Delete user
	return s.uRepo.DeleteUserById(ctx, id)
}

func (s *UserService) checkIfUserExists(ctx context.Context, id int) *e.Error {
	exist, err := s.uRepo.ExistsUserById(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		err = e.NewError(e.LogicUserNotFound, fmt.Sprintf("Could not find user %d.", id))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func (s *UserService) checkIfUsernameIsAlreadyTaken(ctx context.Context, id int,
	username string) *e.Error {
	user, err := s.uRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if (user != nil && id == 0) || (user != nil && user.Id != id) {
		err = e.NewError(e.LogicUserAlreadyExists, fmt.Sprintf("A user with the username '%s' "+
			"already exists.", user.Username))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

// --- User settings functions ---

// GetSettingShowOverviewDetails gets the setting value for the "show overview details" setting.
func (s *UserService) GetSettingShowOverviewDetails(ctx context.Context, userId int) (bool,
	*e.Error) {
	return s.uRepo.GetUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails)
}

// CreateSettingShowOverviewDetails creates the setting value for the "show overview details" setting.
func (s *UserService) CreateSettingShowOverviewDetails(ctx context.Context, userId int,
	value bool) *e.Error {
	return s.uRepo.CreateUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails, value)
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

// CreateUserContract creates the contract information of a user.
func (s *UserService) CreateUserContract(ctx context.Context, userId int,
	contract *model.UserContract) *e.Error {
	return s.uRepo.CreateUserContract(ctx, userId, contract)
}

// UpdateUserContract updates the contract information of a user.
func (s *UserService) UpdateUserContract(ctx context.Context, userId int,
	contract *model.UserContract) *e.Error {
	return s.uRepo.UpdateUserContract(ctx, userId, contract)
}
