package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

// UserRepo retrieves and stores user related entities.
type UserRepo struct {
	repo
}

// NewUserRepo creates a new user repository.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{repo{db}}
}

// --- User functions ---

// GetUsers retrieves all users.
func (r *UserRepo) GetUsers(ctx context.Context) ([]*model.User, error) {
	q := "SELECT id, name, username, password, must_change_password FROM user"

	sr, qErr := r.query(ctx, &scanUserHelper{}, q)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, "Could not query users from database.", qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]*model.User), nil
}

// GetUserById retrieves a user by its ID.
func (r *UserRepo) GetUserById(ctx context.Context, id int) (*model.User, error) {
	q := "SELECT id, name, username, password, must_change_password FROM user WHERE id = ?"

	sr, qErr := r.queryRow(ctx, &scanUserHelper{}, q, id)
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
func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	q := "SELECT id, name, username, password, must_change_password FROM user WHERE username = ?"

	sr, qErr := r.queryRow(ctx, &scanUserHelper{}, q, username)
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
func (r *UserRepo) ExistsUserById(ctx context.Context, id int) (bool, error) {
	cnt, cErr := r.count(ctx, "user", "id = ?", id)
	if cErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user %d from database.",
			id), cErr)
		log.Error(err.StackTrace())
		return false, err
	}

	return cnt > 0, nil
}

// CreateUser creates a new user.
func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	q := "INSERT INTO user (name, username, password, must_change_password) VALUES (?, ?, ?, ?)"

	id, cErr := r.insert(ctx, q, user.Name, user.Username, user.Password, user.MustChangePassword)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, "Could not create user in database.", cErr)
		log.Error(err.StackTrace())
		return err
	}

	user.Id = id

	return nil
}

// UpdateUser updates a user.
func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	q := "UPDATE user SET name = ?, username = ?, password = ?, must_change_password = ? WHERE id = ?"

	uErr := r.exec(ctx, q, user.Name, user.Username, user.Password, user.MustChangePassword, user.Id)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update user %d in database.",
			user.Id), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// DeleteUserById deletes a user by its ID.
func (r *UserRepo) DeleteUserById(ctx context.Context, id int) error {
	q := "DELETE FROM user WHERE id = ?"

	dErr := r.exec(ctx, q, id)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not delete user %d from database.",
			id), dErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

// --- User role functions ---

// GetUserRoles retrieves roles of a user by its ID.
func (r *UserRepo) GetUserRoles(ctx context.Context, userId int) ([]model.Role, error) {
	q := "SELECT r.name FROM user_role ur INNER JOIN role r ON ur.role_id = r.id WHERE ur.user_id = ?"

	sr, qrErr := r.query(ctx, scanRoleHelper{}, q, userId)
	if qrErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not query user roles for user %d "+
			"from database.", userId), qrErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]model.Role), nil
}

// SetUserRoles set roles of a user by its ID.
func (r *UserRepo) SetUserRoles(ctx context.Context, userId int, roles []model.Role) error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) error {
		drErr := r.execWithTx(tx, "DELETE FROM user_role WHERE user_id = ?", userId)
		if drErr != nil {
			err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not update user roles for"+
				" user %d in database.", userId), drErr)
			log.Error(err.StackTrace())
			return err
		}

		rs := make([]string, len(roles))
		for i, role := range roles {
			rs[i] = role.String()
		}

		sel := createSelectionString(rs)
		crErr := r.execWithTx(tx, "INSERT INTO user_role (user_id, role_id) "+
			"SELECT "+strconv.Itoa(userId)+", r.id FROM role r WHERE r.name IN ("+sel+")")
		if crErr != nil {
			err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not update user roles for "+
				"user %d in database.", userId), crErr)
			log.Error(err.StackTrace())
			return err
		}

		return nil
	})
}

// --- User settings functions ---

// GetUserIntSetting retrieves a integer setting of a user.
func (r *UserRepo) GetUserIntSetting(ctx context.Context, userId int, key string) (int, error) {
	v, qErr := r.GetUserStringSetting(ctx, userId, key)
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
func (r *UserRepo) CreateUserIntSetting(ctx context.Context, userId int, key string,
	value int) error {
	return r.CreateUserStringSetting(ctx, userId, key, strconv.Itoa(value))
}

// UpdateUserIntSetting updates a integer setting of a user.
func (r *UserRepo) UpdateUserIntSetting(ctx context.Context, userId int, key string,
	value int) error {
	return r.UpdateUserStringSetting(ctx, userId, key, strconv.Itoa(value))
}

// GetUserBoolSetting retrieves a boolean setting of a user.
func (r *UserRepo) GetUserBoolSetting(ctx context.Context, userId int, key string) (bool, error) {
	v, qErr := r.GetUserStringSetting(ctx, userId, key)
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
func (r *UserRepo) CreateUserBoolSetting(ctx context.Context, userId int, key string, value bool,
) error {
	var v string
	if value {
		v = "true"
	} else {
		v = "false"
	}
	return r.CreateUserStringSetting(ctx, userId, key, v)
}

// UpdateUserBoolSetting updates a boolean setting of a user.
func (r *UserRepo) UpdateUserBoolSetting(ctx context.Context, userId int, key string, value bool,
) error {
	var v string
	if value {
		v = "true"
	} else {
		v = "false"
	}
	return r.UpdateUserStringSetting(ctx, userId, key, v)
}

// GetUserStringSetting retrieves a string setting of a user.
func (r *UserRepo) GetUserStringSetting(ctx context.Context, userId int, key string) (string,
	error) {
	q := "SELECT setting_value FROM user_setting WHERE user_id = ? AND " +
		"setting_key = ?"

	var value string
	qErr := r.queryValue(ctx, &value, q, userId, key)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read user setting '%s' for "+
			"user %d from database.", key, userId), qErr)
		log.Error(err.StackTrace())
		return "", err
	}
	return value, nil
}

// CreateUserStringSetting creates a string setting for a user.
func (r *UserRepo) CreateUserStringSetting(ctx context.Context, userId int, key string,
	value string) error {
	q := "INSERT INTO user_setting (user_id, setting_key, setting_value) VALUES (?, ?, ?)"

	_, cErr := r.insert(ctx, q, userId, key, value)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not create user setting '%s' "+
			"for user %d in database.", key, userId), cErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// UpdateUserStringSetting updates a string setting of a user.
func (r *UserRepo) UpdateUserStringSetting(ctx context.Context, userId int, key string,
	value string) error {
	q := "UPDATE user_setting SET setting_value = ? WHERE user_id = ? AND setting_key = ?"

	uErr := r.exec(ctx, q, value, userId, key)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update user setting '%s' "+
			"for user %d in database.", key, userId), uErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

// --- Helper functions ---

type scanRoleHelper struct {
}

func (h scanRoleHelper) makeSlice() interface{} {
	return make([]model.Role, 0, 10)
}

func (h scanRoleHelper) scan(s scanner) (interface{}, error) {
	var role model.Role

	err := s.Scan(&role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (h scanRoleHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]model.Role), item.(model.Role))
}

type scanUserHelper struct {
}

func (h *scanUserHelper) makeSlice() interface{} {
	return make([]*model.User, 0, 10)
}

func (h *scanUserHelper) scan(s scanner) (interface{}, error) {
	var u model.User

	err := s.Scan(&u.Id, &u.Name, &u.Username, &u.Password, &u.MustChangePassword)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (h *scanUserHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.User), item.(*model.User))
}
