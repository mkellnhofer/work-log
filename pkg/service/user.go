package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/db/repo"
	"kellnhofer.com/work-log/pkg/db/tx"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

// UserService contains user related logic.
type UserService struct {
	service
	uRepo *repo.UserRepo
	cRepo *repo.ContractRepo
}

// NewUserService create a new user service.
func NewUserService(tm *tx.TransactionManager, ur *repo.UserRepo, cr *repo.ContractRepo,
) *UserService {
	return &UserService{service{tm}, ur, cr}
}

// --- Role functions ---

// GetRoles gets all roles.
func (s *UserService) GetRoles(ctx context.Context) ([]model.Role, *e.Error) {
	return model.Roles, nil
}

// GetRolesRights gets all roles with their rights.
func (s *UserService) GetRolesRights(ctx context.Context) (map[model.Role][]model.Right, *e.Error) {
	return model.RolesRights, nil
}

// --- User functions ---

// GetUsers gets all users.
func (s *UserService) GetUsers(ctx context.Context) ([]*model.User, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserData); err != nil {
		return nil, err
	}

	// Get users
	return s.uRepo.GetUsers(ctx)
}

// GetUserById gets a user by its ID.
func (s *UserService) GetUserById(ctx context.Context, id int) (*model.User, *e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, id); err != nil {
		return nil, err
	}

	// Get user
	return s.uRepo.GetUserById(ctx, id)
}

// GetUserByUsername gets a user by its username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User,
	*e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserData); err != nil {
		return nil, err
	}

	// Get user
	return s.uRepo.GetUserByUsername(ctx, username)
}

func (s *UserService) createUser(ctx context.Context, user *model.User) *e.Error {
	// Check if username is already taken
	if err := s.checkIfUsernameIsAlreadyTaken(ctx, 0, user.Username); err != nil {
		return err
	}

	// Hash password
	user.Password = hashUserPassword(user.Password)
	user.MustChangePassword = true

	// Create user
	return s.uRepo.CreateUser(ctx, user)
}

func (s *UserService) updateUser(ctx context.Context, user *model.User) *e.Error {
	// Get user
	oldUser, err := s.getUserById(ctx, user.Id)
	if err != nil {
		return err
	}

	// Check if username is already taken
	if err := s.checkIfUsernameIsAlreadyTaken(ctx, user.Id, user.Username); err != nil {
		return err
	}

	// Was a password provided?
	if user.Password != "" {
		// Use new password
		user.Password = hashUserPassword(user.Password)
		user.MustChangePassword = getCurrentUserId(ctx) != user.Id
	} else {
		// Use old password
		user.Password = oldUser.Password
	}

	// Update user
	return s.uRepo.UpdateUser(ctx, user)
}

// UpdateCurrentUserPassword updates the password of the current a user.
func (s *UserService) UpdateCurrentUserPassword(ctx context.Context, password string) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserAccount); err != nil {
		return err
	}

	// Update user password
	return s.updateUserPassword(ctx, getCurrentUserId(ctx), password)
}

// UpdateUserPassword updates the password of a user.
func (s *UserService) UpdateUserPassword(ctx context.Context, id int, password string) *e.Error {
	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, id); err != nil {
		return err
	}

	// Update user password
	return s.updateUserPassword(ctx, id, password)
}

func (s *UserService) updateUserPassword(ctx context.Context, id int, password string) *e.Error {
	// Get user
	user, err := s.getUserById(ctx, id)
	if err != nil {
		return err
	}

	// Set password
	user.Password = hashUserPassword(password)
	user.MustChangePassword = getCurrentUserId(ctx) != user.Id

	// Update user
	return s.uRepo.UpdateUser(ctx, user)
}

func hashUserPassword(password string) string {
	hBytes, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hErr != nil {
		err := e.WrapError(e.SysUnknown, "Could not hash user password.", hErr)
		log.Error(err.StackTrace())
		panic(err)
	}
	return string(hBytes)
}

// DeleteUserById deletes a user by its ID.
func (s *UserService) DeleteUserById(ctx context.Context, id int) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserData); err != nil {
		return err
	}

	// Check if user exists
	if err := s.checkIfUserExists(ctx, id); err != nil {
		return err
	}

	// Delete user
	return s.uRepo.DeleteUserById(ctx, id)
}

func (s *UserService) getUserById(ctx context.Context, id int) (*model.User, *e.Error) {
	user, err := s.uRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		err = e.NewError(e.LogicUserNotFound, fmt.Sprintf("Could not find user %d.", id))
		log.Debug(err.StackTrace())
		return nil, err
	}
	return user, nil
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

// GetCurrentUserRoles gets the roles of the current user.
func (s *UserService) GetCurrentUserRoles(ctx context.Context) ([]model.Role, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserAccount); err != nil {
		return nil, err
	}

	// Get user roles
	return s.uRepo.GetUserRoles(ctx, getCurrentUserId(ctx))
}

// GetUserRoles gets the roles of a user.
func (s *UserService) GetUserRoles(ctx context.Context, userId int) ([]model.Role, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserData); err != nil {
		return nil, err
	}

	// Get user roles
	return s.uRepo.GetUserRoles(ctx, userId)
}

// SetUserRoles sets the roles of a user.
func (s *UserService) SetUserRoles(ctx context.Context, userId int, roles []model.Role) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserData); err != nil {
		return err
	}

	// Set user roles
	return s.setUserRoles(ctx, userId, roles)
}

func (s *UserService) setUserRoles(ctx context.Context, userId int, roles []model.Role) *e.Error {
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
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return false, err
	}

	// Get setting
	return s.uRepo.GetUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails)
}

func (s *UserService) createSettingShowOverviewDetails(ctx context.Context, userId int,
	value bool) *e.Error {
	return s.uRepo.CreateUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails, value)
}

// UpdateSettingShowOverviewDetails updates the setting value for the "show overview details" setting.
func (s *UserService) UpdateSettingShowOverviewDetails(ctx context.Context, userId int,
	value bool) *e.Error {
	// Check permissions
	if err := s.checkHasCurrentUserChangeRight(ctx, userId); err != nil {
		return err
	}

	// Update setting
	return s.uRepo.UpdateUserBoolSetting(ctx, userId, constant.SettingKeyShowOverviewDetails, value)
}

// --- User contract functions ---

// GetUserContractByUserId gets the contract information of a user by its ID.
func (s *UserService) GetUserContractByUserId(ctx context.Context, userId int) (*model.Contract,
	*e.Error) {
	// Check permissions
	if err := s.checkHasCurrentUserGetRight(ctx, userId); err != nil {
		return nil, err
	}

	// Get user contract
	return s.cRepo.GetContractByUserId(ctx, userId)
}

func (s *UserService) createUserContract(ctx context.Context, userId int,
	contract *model.Contract) *e.Error {
	// Check contract
	if err := s.checkUserContract(contract); err != nil {
		return err
	}

	// Create contract
	return s.cRepo.CreateContract(ctx, userId, contract)
}

func (s *UserService) updateUserContract(ctx context.Context, userId int,
	contract *model.Contract) *e.Error {
	// Check contract
	if err := s.checkUserContract(contract); err != nil {
		return err
	}

	// Update contract
	return s.cRepo.UpdateContract(ctx, userId, contract)
}

func (s *UserService) checkUserContract(contract *model.Contract) *e.Error {
	if err := s.checkUserContractWorkingHours(contract.FirstDay, contract.WorkingHours); err != nil {
		return err
	}
	if err := s.checkUserContractVacationDays(contract.FirstDay, contract.VacationDays); err != nil {
		return err
	}
	return nil
}

func (s *UserService) checkUserContractWorkingHours(contractFirstDay time.Time,
	workingHours []model.ContractWorkingHours) *e.Error {
	errCode := e.LogicContractWorkingHoursInvalid

	// Check if intervals are empty
	if len(workingHours) == 0 {
		err := e.NewError(errCode, "Working hours intervals cannot be empty.")
		log.Debug(err.StackTrace())
		return err
	}

	// Check intervals
	startDay := time.Time{}
	for _, whs := range workingHours {
		// Get start day of intervals
		if startDay.IsZero() || whs.FirstDay.Before(startDay) {
			startDay = whs.FirstDay
		}

		// Check if interval first day is not first day of month
		if whs.FirstDay.Day() != 1 {
			err := e.NewError(errCode, "A working hours interval must start at the first day of a "+
				"month.")
			log.Debug(err.StackTrace())
			return err
		}
		// Check if interval hours are negative
		if whs.Hours < 0 {
			err := e.NewError(errCode, "Working hours cannot be negative.")
			log.Debug(err.StackTrace())
			return err
		}
	}

	// Check if start day does not matches with contract first day
	if startDay.After(contractFirstDay) {
		err := e.NewError(errCode, "Working hours intervals start must match contract start.")
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

func (s *UserService) checkUserContractVacationDays(contractFirstDay time.Time,
	vacationDays []model.ContractVacationDays) *e.Error {
	errCode := e.LogicContractVacationDaysInvalid

	// Check if intervals are empty
	if len(vacationDays) == 0 {
		err := e.NewError(errCode, "Vacation days intervals cannot be empty.")
		log.Debug(err.StackTrace())
		return err
	}

	// Check intervals
	startDay := time.Time{}
	for _, vds := range vacationDays {
		// Get start day of intervals
		if startDay.IsZero() || vds.FirstDay.Before(startDay) {
			startDay = vds.FirstDay
		}

		// Check if interval first day is not first day of month
		if vds.FirstDay.Day() != 1 {
			err := e.NewError(errCode, "A vacation days interval must start at the first day of a "+
				"month.")
			log.Debug(err.StackTrace())
			return err
		}
		// Check if interval days are negative
		if vds.Days < 0 {
			err := e.NewError(errCode, "Vacation days cannot be negative.")
			log.Debug(err.StackTrace())
			return err
		}
	}

	// Check if start day does not matches with contract first day
	if startDay.After(contractFirstDay) {
		err := e.NewError(errCode, "Vacation days intervals start must match contract start.")
		log.Debug(err.StackTrace())
		return err
	}

	return nil
}

// --- User data functions ---

// GetUserDatas gets all users with related information at once.
func (s *UserService) GetUserDatas(ctx context.Context) ([]*model.UserData, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserData); err != nil {
		return nil, err
	}

	// Get users
	users, gusErr := s.uRepo.GetUsers(ctx)
	if gusErr != nil {
		return nil, gusErr
	}

	userDatas := make([]*model.UserData, 0, 10)

	// Get user contracts
	for _, user := range users {
		contract, gucErr := s.cRepo.GetContractByUserId(ctx, user.Id)
		if gucErr != nil {
			return nil, gucErr
		}

		userDatas = append(userDatas, model.NewUserData(user.Id, user, contract))
	}

	return userDatas, nil
}

// GetCurrentUserData gets the current user with related information at once.
func (s *UserService) GetCurrentUserData(ctx context.Context) (*model.UserData, *e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserAccount); err != nil {
		return nil, err
	}

	userId := getCurrentUserId(ctx)

	// Get user
	user, guErr := s.uRepo.GetUserById(ctx, userId)
	if guErr != nil {
		return nil, guErr
	}
	if user == nil {
		return nil, nil
	}

	// Get user contract
	contract, gucErr := s.cRepo.GetContractByUserId(ctx, userId)
	if gucErr != nil {
		return nil, gucErr
	}

	return model.NewUserData(userId, user, contract), nil
}

// GetUserDataByUserId gets a user with related information at once.
func (s *UserService) GetUserDataByUserId(ctx context.Context, userId int) (*model.UserData,
	*e.Error) {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightGetUserData); err != nil {
		return nil, err
	}

	// Get user
	user, guErr := s.uRepo.GetUserById(ctx, userId)
	if guErr != nil {
		return nil, guErr
	}
	if user == nil {
		return nil, nil
	}

	// Get user contract
	contract, gucErr := s.cRepo.GetContractByUserId(ctx, userId)
	if gucErr != nil {
		return nil, gucErr
	}

	return model.NewUserData(userId, user, contract), nil
}

// CreateUserData creates a new user with related information at once.
func (s *UserService) CreateUserData(ctx context.Context, userData *model.UserData) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserData); err != nil {
		return err
	}

	// Execute in transaction
	return s.tm.ExecuteInNewTransaction(ctx, func(ctx context.Context) *e.Error {
		// Create user
		if err := s.createUser(ctx, userData.User); err != nil {
			return err
		}
		userData.Id = userData.User.Id
		// Create roles
		userRoles := []model.Role{model.RoleUser}
		if err := s.setUserRoles(ctx, userData.Id, userRoles); err != nil {
			return err
		}
		// Create settings
		if err := s.createSettingShowOverviewDetails(ctx, userData.Id, true); err != nil {
			return err
		}
		// Create contract
		return s.createUserContract(ctx, userData.Id, userData.Contract)
	})
}

// UpdateUserData updates a user with related information at once.
func (s *UserService) UpdateUserData(ctx context.Context, userData *model.UserData) *e.Error {
	// Check permissions
	if err := checkHasCurrentUserRight(ctx, model.RightChangeUserData); err != nil {
		return err
	}

	// Execute in transaction
	return s.tm.ExecuteInNewTransaction(ctx, func(ctx context.Context) *e.Error {
		// Update user
		if userData.User != nil {
			if err := s.updateUser(ctx, userData.User); err != nil {
				return err
			}
		}
		// Update contract
		if userData.Contract != nil {
			if err := s.updateUserContract(ctx, userData.Id, userData.Contract); err != nil {
				return err
			}
		}
		return nil
	})
}

// --- Permission helper functions ---

func (s *UserService) checkHasCurrentUserGetRight(ctx context.Context, userId int) *e.Error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightGetUserAccount)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightGetUserData)
	}
}

func (s *UserService) checkHasCurrentUserChangeRight(ctx context.Context, userId int) *e.Error {
	if userId == getCurrentUserId(ctx) {
		return checkHasCurrentUserRight(ctx, model.RightChangeUserAccount)
	} else {
		return checkHasCurrentUserRight(ctx, model.RightChangeUserData)
	}
}
