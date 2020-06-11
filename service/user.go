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

// --- Role functions ---

// GetRoles gets all roles.
func (s *UserService) GetRoles(ctx context.Context) []model.Role {
	return model.Roles
}

// GetRolesRights gets all roles with their rights.
func (s *UserService) GetRolesRights(ctx context.Context) map[model.Role][]model.Right {
	return model.RolesRights
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

// --- User role functions ---

// GetUserRoles gets the roles of a user.
func (s *UserService) GetUserRoles(ctx context.Context, userId int) ([]model.Role, *e.Error) {
	return s.uRepo.GetUserRoles(ctx, userId)
}

// SetUserRoles sets the roles of a user.
func (s *UserService) SetUserRoles(ctx context.Context, userId int, roles []model.Role) *e.Error {
	// Check if roles exist
	if err := s.checkIfRolesExist(ctx, roles); err != nil {
		return err
	}

	// Set user roles
	return s.uRepo.SetUserRoles(ctx, userId, roles)
}

func (s *UserService) checkIfRolesExist(ctx context.Context, roles []model.Role) *e.Error {
	for _, role := range roles {
		found := containsRole(model.Roles, role)
		if !found {
			err := e.NewError(e.LogicRoleNotFound, fmt.Sprintf("Role '%s' does not exists.", role))
			log.Debug(err.StackTrace())
			return err
		}
	}
	return nil
}

func containsRole(roles []model.Role, role model.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
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

// --- User data functions ---

// GetUserDatas gets all users with related information at once.
func (s *UserService) GetUserDatas(ctx context.Context) ([]*model.UserData, *e.Error) {
	users, gusErr := s.GetUsers(ctx)
	if gusErr != nil {
		return nil, gusErr
	}

	userDatas := make([]*model.UserData, 0, 10)

	for _, user := range users {
		userContract, gucErr := s.GetUserContractByUserId(ctx, user.Id)
		if gucErr != nil {
			return nil, gucErr
		}

		userDatas = append(userDatas, model.NewUserData(user.Id, user, userContract))
	}

	return userDatas, nil
}

// GetUserDataByUserId gets a user with related information at once.
func (s *UserService) GetUserDataByUserId(ctx context.Context, userId int) (*model.UserData,
	*e.Error) {
	user, guErr := s.GetUserById(ctx, userId)
	if guErr != nil {
		return nil, guErr
	}
	if user == nil {
		return nil, nil
	}

	userContract, gucErr := s.GetUserContractByUserId(ctx, userId)
	if gucErr != nil {
		return nil, gucErr
	}

	return model.NewUserData(userId, user, userContract), nil
}

// CreateUserData creates a new user with related information at once.
func (s *UserService) CreateUserData(ctx context.Context, userData *model.UserData) *e.Error {
	// Start transaction
	if err := s.tm.Begin(ctx); err != nil {
		return err
	}

	// Create user
	if err := s.CreateUser(ctx, userData.User); err != nil {
		s.tm.Rollback(ctx)
		return err
	}
	userData.Id = userData.User.Id
	// Create roles
	userRoles := []model.Role{model.RoleUser}
	if err := s.SetUserRoles(ctx, userData.Id, userRoles); err != nil {
		s.tm.Rollback(ctx)
		return err
	}
	// Create settings
	if err := s.CreateSettingShowOverviewDetails(ctx, userData.Id, true); err != nil {
		s.tm.Rollback(ctx)
		return err
	}
	// Create contract
	if err := s.CreateUserContract(ctx, userData.Id, userData.UserContract); err != nil {
		s.tm.Rollback(ctx)
		return err
	}

	// Commit transaction
	return s.tm.Commit(ctx)
}

// UpdateUserData updates a user with related information at once.
func (s *UserService) UpdateUserData(ctx context.Context, userData *model.UserData) *e.Error {
	// Start transaction
	if err := s.tm.Begin(ctx); err != nil {
		return err
	}

	// Update user
	if userData.User != nil {
		if err := s.UpdateUser(ctx, userData.User); err != nil {
			s.tm.Rollback(ctx)
			return err
		}
	}
	// Update contract
	if userData.UserContract != nil {
		if err := s.UpdateUserContract(ctx, userData.Id, userData.UserContract); err != nil {
			s.tm.Rollback(ctx)
			return err
		}
	}

	// Commit transaction
	return s.tm.Commit(ctx)
}
