package mapper

import (
	"kellnhofer.com/work-log/pkg/model"
	vm "kellnhofer.com/work-log/web/model"
)

// UserMapper creates view models for the user page.
type UserMapper struct {
	mapper
}

// NewUserMapper creates a new user mapper.
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// CreateUserInfoViewModel creates a view model for basic user information.
func (m *mapper) CreateUserInfoViewModel(user *model.User) *vm.UserInfo {
	return &vm.UserInfo{
		Id:       user.Id,
		Initials: getUserInitials(user.Name),
	}
}
