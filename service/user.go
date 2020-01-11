package service

import (
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
