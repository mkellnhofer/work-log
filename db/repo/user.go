package repo

import (
	"database/sql"
	"fmt"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

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
