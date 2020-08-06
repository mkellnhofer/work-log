package validator

import (
	"fmt"
	"regexp"

	vm "kellnhofer.com/work-log/api/model"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	m "kellnhofer.com/work-log/model"
)

// --- User API model valdidation functions ---

// ValidateCreateUser validates information of a CreateUserData API model.
func ValidateCreateUser(data *vm.CreateUserData) *e.Error {
	if err := checkUserName(data.Name); err != nil {
		return err
	}
	if err := checkUserUsername(data.Username); err != nil {
		return err
	}
	if err := checkUserPassword(data.Password); err != nil {
		return err
	}
	return ValidateCreateUserContract(data.Contract)
}

// ValidateCreateUserContract validates information of a CreateUserContract API model.
func ValidateCreateUserContract(data *vm.CreateUserContract) *e.Error {
	if err := checkNotNil("contract", data); err != nil {
		return err
	}
	if err := checkUserDailyWorkingDuration(data.DailyWorkingDuration); err != nil {
		return err
	}
	if err := checkUserAnnualVacationDays(data.AnnualVacationDays); err != nil {
		return err
	}
	if err := checkUserInitOvertimeDuration(data.InitOvertimeDuration); err != nil {
		return err
	}
	if err := checkUserInitVacationDays(data.InitVacationDays); err != nil {
		return err
	}
	return checkUserFirstWorkDay(data.FirstWorkDay)
}

// ValidateUpdateUser validates information of a UpdateUserData API model.
func ValidateUpdateUser(data *vm.UpdateUserData) *e.Error {
	if err := checkUserName(data.Name); err != nil {
		return err
	}
	if err := checkUserUsername(data.Username); err != nil {
		return err
	}
	return ValidateUpdateUserContract(data.Contract)
}

// ValidateUpdateUserContract validates information of a UpdateUserContract API model.
func ValidateUpdateUserContract(data *vm.UpdateUserContract) *e.Error {
	if err := checkNotNil("contract", data); err != nil {
		return err
	}
	if err := checkUserDailyWorkingDuration(data.DailyWorkingDuration); err != nil {
		return err
	}
	if err := checkUserAnnualVacationDays(data.AnnualVacationDays); err != nil {
		return err
	}
	if err := checkUserInitOvertimeDuration(data.InitOvertimeDuration); err != nil {
		return err
	}
	if err := checkUserInitVacationDays(data.InitVacationDays); err != nil {
		return err
	}
	return checkUserFirstWorkDay(data.FirstWorkDay)
}

// ValidateUpdateUserPassword validates information of a UpdateUserPassword API model.
func ValidateUpdateUserPassword(data *vm.UpdateUserPassword) *e.Error {
	return checkUserPassword(data.Password)
}

// ValidateUpdateUserRoles validates information of a UpdateUserRoles API model.
func ValidateUpdateUserRoles(data *vm.UpdateUserRoles) *e.Error {
	if err := checkArrayLengthNotZero("roles", len(data.Roles)); err != nil {
		return err
	}
	if err := checkStringArrayNotEmpty("roles", data.Roles); err != nil {
		return err
	}
	if err := checkStringArrayNotTooLong("roles", data.Roles, m.MaxLengthRoleName); err != nil {
		return err
	}
	for _, role := range data.Roles {
		if err := checkRole(role); err != nil {
			return err
		}
	}
	return nil
}

// --- Basic user validation functions ---

func checkRole(role string) *e.Error {
	r := regexp.MustCompile("^[a-z_]+$")
	if !r.MatchString(role) {
		err := e.NewError(e.ValRoleInvalid, fmt.Sprintf("'role' must only contain letters and "+
			"following special characters '%s'.", "_"))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkUserName(name string) *e.Error {
	if err := checkStringNotEmpty("name", name); err != nil {
		return err
	}
	if err := checkStringNotTooLong("name", name, m.MaxLengthUserName); err != nil {
		return err
	}
	return nil
}

func checkUserUsername(username string) *e.Error {
	if len(username) == 0 {
		err := e.NewError(e.ValUsernameInvalid, "'username' must not be empty.")
		log.Debug(err.StackTrace())
		return err
	}
	if len(username) < 4 {
		err := e.NewError(e.ValUsernameInvalid, fmt.Sprintf("'username' must be at least %d long.",
			4))
		log.Debug(err.StackTrace())
		return err
	}
	if len(username) > m.MaxLengthUserUsername {
		err := e.NewError(e.ValUsernameInvalid, fmt.Sprintf("'username' must not be longer than %d.",
			m.MaxLengthUserUsername))
		log.Debug(err.StackTrace())
		return err
	}
	r := regexp.MustCompile("^[0-9a-zA-Z\\-.]+$")
	if !r.MatchString(username) {
		err := e.NewError(e.ValUsernameInvalid, "'username' contains contains illegal character.")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkUserPassword(password string) *e.Error {
	if len(password) == 0 {
		err := e.NewError(e.ValPasswordInvalid, "'password' must not be empty.")
		log.Debug(err.StackTrace())
		return err
	}
	if len(password) < 8 {
		err := e.NewError(e.ValPasswordInvalid, fmt.Sprintf("'password' must be at least %d long.",
			8))
		log.Debug(err.StackTrace())
		return err
	}
	if len(password) > m.MaxLengthUserPassword {
		err := e.NewError(e.ValPasswordInvalid, fmt.Sprintf("'password' must not be longer than %d.",
			m.MaxLengthUserPassword))
		log.Debug(err.StackTrace())
		return err
	}
	r := regexp.MustCompile("^[0-9a-zA-Z!\"#$%&'()*+,\\-./:;=?@\\[\\\\\\]^_{|}~]+$")
	if !r.MatchString(password) {
		err := e.NewError(e.ValPasswordInvalid, "'password' contains contains illegal character.")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkUserDailyWorkingDuration(num float32) *e.Error {
	return checkFloatNotNegativeOrZero("dailyWorkingDuration", num)
}

func checkUserAnnualVacationDays(num float32) *e.Error {
	return checkFloatNotNegative("annualVacationDays", num)
}

func checkUserInitOvertimeDuration(num float32) *e.Error {
	return nil
}

func checkUserInitVacationDays(num float32) *e.Error {
	return nil
}

func checkUserFirstWorkDay(date string) *e.Error {
	return checkDateValid("firstWorkDay", date)
}
