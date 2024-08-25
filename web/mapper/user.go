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

// CreateUserProfileInfoViewModel creates a view model for detailed user information.
func (m *UserMapper) CreateUserProfileInfoViewModel(user *model.User, contract *model.Contract,
	) *vm.UserProfileInfo {
	profileInfo := &vm.UserProfileInfo{
		Id:       user.Id,
		Initials: getUserInitials(user.Name),
		Name:     user.Name,
		Username: user.Username,
	}
	if contract != nil {
		ci := &vm.ContractInfo{
			FirstDay:          formatDate(contract.FirstDay),
			InitOvertimeHours: getHoursString(contract.InitOvertimeHours),
			InitVacationDays:  getDaysString(contract.InitVacationDays),
		}
		for _, wh := range contract.WorkingHours {
			ci.WorkingHours = append(ci.WorkingHours, &vm.ContractWorkingHours{
				FirstDay: formatDate(wh.FirstDay),
				Hours:    getHoursString(wh.Hours),
			})
		}
		for _, vd := range contract.VacationDays {
			ci.VacationDays = append(ci.VacationDays, &vm.ContractVacationDays{
				FirstDay: formatDate(vd.FirstDay),
				Days:     getDaysString(vd.Days),
			})
		}
		profileInfo.Contract = ci
	}
	return profileInfo
}
