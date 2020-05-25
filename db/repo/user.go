package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

type dbUserContract struct {
	dailyWorkingDuration int
	annualVacationDays   float32
	initOvertimeDuration int
	initVacationDays     float32
	firstWorkDay         string
}

// UserRepo retrieves and stores user and role records.
type UserRepo struct {
	repo
}

// NewUserRepo creates a new user repository.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{repo{db}}
}

// --- User functions ---

// GetUsers retrieves all users.
func (r *UserRepo) GetUsers() ([]*model.User, *e.Error) {
	q := "SELECT id, name, username, password FROM user"

	sr, qErr := r.query(&scanUserHelper{}, q)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query users from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]*model.User), nil
}

// GetUserById retrieves a user by its ID.
func (r *UserRepo) GetUserById(id int) (*model.User, *e.Error) {
	q := "SELECT id, name, username, password FROM user WHERE id = ?"

	sr, qErr := r.queryRow(&scanUserHelper{}, q, id)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user %d from database.",
			id), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.User), nil
}

// GetUserByUsername retrieves a user by its username.
func (r *UserRepo) GetUserByUsername(username string) (*model.User, *e.Error) {
	q := "SELECT id, name, username, password FROM user WHERE username = ?"

	sr, qErr := r.queryRow(&scanUserHelper{}, q, username)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user '%s' from database.",
			username), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.User), nil
}

// ExistsUserById checks if a user exists.
func (r *UserRepo) ExistsUserById(id int) (bool, *e.Error) {
	cnt, cErr := r.count("user", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user %d from database.",
			id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateUser creates a new user.
func (r *UserRepo) CreateUser(user *model.User) *e.Error {
	q := "INSERT INTO user (name, username, password) VALUES (?, ?, ?)"

	id, cErr := r.insert(q, user.Name, user.Username, user.Password)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create user in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	user.Id = id

	return nil
}

// UpdateUser updates a user.
func (r *UserRepo) UpdateUser(user *model.User) *e.Error {
	q := "UPDATE user SET name = ?, username = ?, password = ? WHERE id = ?"

	uErr := r.exec(q, user.Name, user.Username, user.Password, user.Id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update user %d in database.",
			user.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteUserById deletes a user by its ID.
func (r *UserRepo) DeleteUserById(id int) *e.Error {
	q := "DELETE FROM user WHERE id = ?"

	dErr := r.exec(q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete user %d from database.",
			id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- User settings functions ---

// GetUserIntSetting retrieves a integer setting of a user.
func (r *UserRepo) GetUserIntSetting(userId int, key string) (int, *e.Error) {
	v, qErr := r.GetUserStringSetting(userId, key)
	if qErr != nil {
		return 0, qErr
	}

	value, cErr := strconv.Atoi(v)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user setting '%s' for "+
			"user %d from database.", key, userId), cErr)
		log.Error(err.StackTrace())
		return 0, err
	}
	return value, nil
}

// CreateUserIntSetting creates a integer setting for a user.
func (r *UserRepo) CreateUserIntSetting(userId int, key string, value int) *e.Error {
	return r.CreateUserStringSetting(userId, key, strconv.Itoa(value))
}

// UpdateUserIntSetting updates a integer setting of a user.
func (r *UserRepo) UpdateUserIntSetting(userId int, key string, value int) *e.Error {
	return r.UpdateUserStringSetting(userId, key, strconv.Itoa(value))
}

// GetUserBoolSetting retrieves a boolean setting of a user.
func (r *UserRepo) GetUserBoolSetting(userId int, key string) (bool, *e.Error) {
	v, qErr := r.GetUserStringSetting(userId, key)
	if qErr != nil {
		return false, qErr
	}
	if v == "true" {
		return true, nil
	} else if v == "false" {
		return false, nil
	} else {
		err := e.NewError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user setting '%s' for "+
			"user %d from database.", key, userId))
		log.Error(err.StackTrace())
		return false, err
	}
}

// CreateUserBoolSetting creates a boolean setting for a user.
func (r *UserRepo) CreateUserBoolSetting(userId int, key string, value bool) *e.Error {
	var v string
	if value {
		v = "true"
	} else {
		v = "false"
	}
	return r.CreateUserStringSetting(userId, key, v)
}

// UpdateUserBoolSetting updates a boolean setting of a user.
func (r *UserRepo) UpdateUserBoolSetting(userId int, key string, value bool) *e.Error {
	var v string
	if value {
		v = "true"
	} else {
		v = "false"
	}
	return r.UpdateUserStringSetting(userId, key, v)
}

// GetUserStringSetting retrieves a string setting of a user.
func (r *UserRepo) GetUserStringSetting(userId int, key string) (string, *e.Error) {
	q := "SELECT setting_value FROM user_setting WHERE user_id = ? AND " +
		"setting_key = ?"

	var value string
	qErr := r.queryValue(&value, q, userId, key)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user setting '%s' for "+
			"user %d from database.", key, userId), qErr)
		log.Error(err.StackTrace())
		return "", err
	}
	return value, nil
}

// CreateUserStringSetting creates a string setting for a user.
func (r *UserRepo) CreateUserStringSetting(userId int, key string, value string) *e.Error {
	q := "INSERT INTO user_setting (user_id, setting_key, setting_value) VALUES (?, ?, ?)"

	_, cErr := r.insert(q, userId, key, value)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not create user setting '%s' "+
			"for user %d in database.", key, userId), cErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// UpdateUserStringSetting updates a string setting of a user.
func (r *UserRepo) UpdateUserStringSetting(userId int, key string, value string) *e.Error {
	q := "UPDATE user_setting SET setting_value = ? WHERE user_id = ? AND setting_key = ?"

	uErr := r.exec(q, value, userId, key)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update user setting '%s' "+
			"for user %d in database.", key, userId), uErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// --- User contract functions ---

// GetUserContractByUserId retrieves the contract information of a user by its ID.
func (r *UserRepo) GetUserContractByUserId(userId int) (*model.UserContract, *e.Error) {
	q := "SELECT daily_working_duration, annual_vacation_days, init_overtime_duration, " +
		"init_vacation_days, first_work_day FROM user_contract WHERE user_id = ?"

	sr, qErr := r.queryRow(&scanUserContractHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user contract for user "+
			"%d from database.", userId), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	if sr == nil {
		return nil, nil
	}
	return sr.(*model.UserContract), nil
}

// CreateUserContract creates the contract information of a user.
func (r *UserRepo) CreateUserContract(userId int, userContract *model.UserContract) *e.Error {
	uc := toDbUserContract(userContract)

	q := "INSERT INTO user_contract (user_id, daily_working_duration, annual_vacation_days, " +
		"init_overtime_duration, init_vacation_days, first_work_day) VALUES (?, ?, ?, ?, ?, ?)"

	_, cErr := r.insert(q, userId, uc.dailyWorkingDuration, uc.annualVacationDays,
		uc.initOvertimeDuration, uc.initVacationDays, uc.firstWorkDay)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not create user contract for "+
			"user %d from database.", userId), cErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// UpdateUserContract updates the contract information of a user.
func (r *UserRepo) UpdateUserContract(userId int, userContract *model.UserContract) *e.Error {
	uc := toDbUserContract(userContract)

	q := "UPDATE user_contract SET daily_working_duration = ?, annual_vacation_days = ?, " +
		"init_overtime_duration = ?, init_vacation_days = ?, first_work_day = ? WHERE user_id = ?"

	uErr := r.exec(q, uc.dailyWorkingDuration, uc.annualVacationDays, uc.initOvertimeDuration,
		uc.initVacationDays, uc.firstWorkDay, userId)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update user contract for "+
			"user %d in database.", userId), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- Helper functions ---

type scanUserHelper struct {
}

func (h *scanUserHelper) makeSlice() interface{} {
	return make([]*model.User, 0, 10)
}

func (h *scanUserHelper) scan(s scanner) (interface{}, error) {
	var u model.User

	err := s.Scan(&u.Id, &u.Name, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (h *scanUserHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.User), item.(*model.User))
}

type scanUserContractHelper struct {
}

func (h *scanUserContractHelper) makeSlice() interface{} {
	return make([]*model.UserContract, 0, 10)
}

func (h *scanUserContractHelper) scan(s scanner) (interface{}, error) {
	var dbUc dbUserContract

	err := s.Scan(&dbUc.dailyWorkingDuration, &dbUc.annualVacationDays, &dbUc.initOvertimeDuration,
		&dbUc.initVacationDays, &dbUc.firstWorkDay)
	if err != nil {
		return nil, err
	}

	uc := fromDbUserContract(&dbUc)

	return uc, nil
}

func (h *scanUserContractHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.UserContract), item.(*model.UserContract))
}

func toDbUserContract(in *model.UserContract) *dbUserContract {
	var out dbUserContract
	out.dailyWorkingDuration = *formatDuration(&in.DailyWorkingDuration)
	out.annualVacationDays = in.AnnualVacationDays
	out.initOvertimeDuration = *formatDuration(&in.InitOvertimeDuration)
	out.initVacationDays = in.InitVacationDays
	out.firstWorkDay = *formatTimestamp(&in.FirstWorkDay)
	return &out
}

func fromDbUserContract(in *dbUserContract) *model.UserContract {
	var out model.UserContract
	out.DailyWorkingDuration = *parseDuration(&in.dailyWorkingDuration)
	out.AnnualVacationDays = in.annualVacationDays
	out.InitOvertimeDuration = *parseDuration(&in.initOvertimeDuration)
	out.InitVacationDays = in.initVacationDays
	out.FirstWorkDay = *parseTimestamp(&in.firstWorkDay)
	return &out
}
