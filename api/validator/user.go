package validator

import (
	"fmt"
	"regexp"

	vm "kellnhofer.com/work-log/api/model"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	m "kellnhofer.com/work-log/pkg/model"
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
	return ValidateCreateContract(data.Contract)
}

// ValidateCreateContract validates information of a CreateContract API model.
func ValidateCreateContract(data *vm.CreateContract) *e.Error {
	if err := checkNotNil("contract", data); err != nil {
		return err
	}
	if err := checkContractFirstDay(data.FirstDay); err != nil {
		return err
	}
	if err := checkContractInitOvertimeHours(data.InitOvertimeHours); err != nil {
		return err
	}
	if err := checkContractInitVacationDays(data.InitVacationDays); err != nil {
		return err
	}
	if err := checkContractWorkingHours(data.WorkingHours); err != nil {
		return err
	}
	return checkContractVacationDays(data.VacationDays)
}

// ValidateUpdateUser validates information of a UpdateUserData API model.
func ValidateUpdateUser(data *vm.UpdateUserData) *e.Error {
	if err := checkUserName(data.Name); err != nil {
		return err
	}
	if err := checkUserUsername(data.Username); err != nil {
		return err
	}
	return ValidateUpdateContract(data.Contract)
}

// ValidateUpdateContract validates information of a UpdateContract API model.
func ValidateUpdateContract(data *vm.UpdateContract) *e.Error {
	if err := checkNotNil("contract", data); err != nil {
		return err
	}
	if err := checkContractFirstDay(data.FirstDay); err != nil {
		return err
	}
	if err := checkContractInitOvertimeHours(data.InitOvertimeHours); err != nil {
		return err
	}
	if err := checkContractInitVacationDays(data.InitVacationDays); err != nil {
		return err
	}
	if err := checkContractWorkingHours(data.WorkingHours); err != nil {
		return err
	}
	return checkContractVacationDays(data.VacationDays)
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
	if len(username) < m.MinLengthUserUsername {
		err := e.NewError(e.ValUsernameInvalid, fmt.Sprintf("'username' must be at least %d long.",
			m.MinLengthUserUsername))
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
	if len(password) < m.MinLengthUserPassword {
		err := e.NewError(e.ValPasswordInvalid, fmt.Sprintf("'password' must be at least %d long.",
			m.MinLengthUserPassword))
		log.Debug(err.StackTrace())
		return err
	}
	if len(password) > m.MaxLengthUserPassword {
		err := e.NewError(e.ValPasswordInvalid, fmt.Sprintf("'password' must not be longer than %d.",
			m.MaxLengthUserPassword))
		log.Debug(err.StackTrace())
		return err
	}
	r := regexp.MustCompile("^[" + m.ValidUserPasswordCharacters + "]+$")
	if !r.MatchString(password) {
		err := e.NewError(e.ValPasswordInvalid, "'password' contains contains illegal character.")
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkContractFirstDay(date string) *e.Error {
	return checkDateValid("firstDay", date)
}

func checkContractInitOvertimeHours(num float32) *e.Error {
	// Nothing to validate, hours can be <0, 0 and >0
	return nil
}

func checkContractInitVacationDays(num float32) *e.Error {
	// Nothing to validate, days can be <0, 0 and >0
	return nil
}

func checkContractWorkingHours(data []*vm.ContractWorkingHours) *e.Error {
	if err := checkArrayLengthNotZero("workingHours", len(data)); err != nil {
		return err
	}
	for _, wh := range data {
		if wh == nil {
			err := e.NewError(e.ValFieldNil, "Elements of 'workingHours' must not be null.")
			log.Debug(err.StackTrace())
			return err
		}
		if err := checkContractDailyWorkingHours(wh); err != nil {
			return err
		}
	}
	return nil
}

func checkContractDailyWorkingHours(data *vm.ContractWorkingHours) *e.Error {
	if err := checkDateValid("firstDay", data.FirstDay); err != nil {
		return err
	}
	return checkFloatNotNegativeOrZero("hours", data.Hours)
}

func checkContractVacationDays(data []*vm.ContractVacationDays) *e.Error {
	if err := checkArrayLengthNotZero("vacationDays", len(data)); err != nil {
		return err
	}
	for _, vd := range data {
		if vd == nil {
			err := e.NewError(e.ValFieldNil, "Elements of 'vacationDays' must not be null.")
			log.Debug(err.StackTrace())
			return err
		}
		if err := checkContractMonthlyVacationDays(vd); err != nil {
			return err
		}
	}
	return nil
}

func checkContractMonthlyVacationDays(data *vm.ContractVacationDays) *e.Error {
	if err := checkDateValid("firstDay", data.FirstDay); err != nil {
		return err
	}
	return checkFloatNotNegative("days", data.Days)
}
