package mapper

import (
	"strings"

	am "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/model"
)

// --- User functions ---

// ToUserDatas converts a list of logic user models to a list of API user models.
func ToUserDatas(uds []*m.UserData) *am.UserDataList {
	if uds == nil {
		return nil
	}

	items := make([]*am.UserData, len(uds))
	for i, ud := range uds {
		items[i] = ToUserData(ud)
	}

	return am.NewUserDataList(items)
}

// ToUserData converts a logic user model to an API user model.
func ToUserData(ud *m.UserData) *am.UserData {
	if ud == nil {
		return nil
	}

	var out am.UserData
	out.Id = ud.Id
	out.Name = ud.User.Name
	out.Username = ud.User.Username
	out.Contract = toUserContract(ud.UserContract)
	return &out
}

// FromCreateUserData converts an API user creation model to a logic user data model.
func FromCreateUserData(cud *am.CreateUserData) *m.UserData {
	if cud == nil {
		return nil
	}

	var out m.UserData
	out.User = fromCreateUser(cud)
	out.UserContract = fromCreateUserContract(cud.Contract)
	return &out
}

// FromUpdateUserData converts an API user update model to a logic user data model.
func FromUpdateUserData(id int, uud *am.UpdateUserData) *m.UserData {
	if uud == nil {
		return nil
	}

	var out m.UserData
	out.Id = id
	out.User = fromUpdateUser(id, uud)
	out.UserContract = fromUpdateUserContract(uud.Contract)
	return &out
}

func fromCreateUser(cud *am.CreateUserData) *m.User {
	if cud == nil {
		return nil
	}

	var out m.User
	out.Name = strings.TrimSpace(cud.Name)
	out.Username = strings.TrimSpace(cud.Username)
	out.Password = strings.TrimSpace(cud.Password)
	return &out
}

func fromUpdateUser(id int, uud *am.UpdateUserData) *m.User {
	if uud == nil {
		return nil
	}

	var out m.User
	out.Id = id
	out.Name = strings.TrimSpace(uud.Name)
	out.Username = strings.TrimSpace(uud.Username)
	return &out
}

func toUserContract(uc *m.UserContract) *am.UserContract {
	if uc == nil {
		return nil
	}

	var out am.UserContract
	out.DailyWorkingDuration = formatHoursDuration(uc.DailyWorkingDuration)
	out.AnnualVacationDays = uc.AnnualVacationDays
	out.InitOvertimeDuration = formatHoursDuration(uc.InitOvertimeDuration)
	out.InitVacationDays = uc.InitVacationDays
	out.FirstWorkDay = formatDate(uc.FirstWorkDay)
	return &out
}

func fromCreateUserContract(cuc *am.CreateUserContract) *m.UserContract {
	if cuc == nil {
		return nil
	}

	var out m.UserContract
	out.DailyWorkingDuration = parseHoursDuration(cuc.DailyWorkingDuration)
	out.AnnualVacationDays = cuc.AnnualVacationDays
	out.InitOvertimeDuration = parseHoursDuration(cuc.InitOvertimeDuration)
	out.InitVacationDays = cuc.InitVacationDays
	out.FirstWorkDay = parseDate(cuc.FirstWorkDay)
	return &out
}

func fromUpdateUserContract(uuc *am.UpdateUserContract) *m.UserContract {
	if uuc == nil {
		return nil
	}

	var out m.UserContract
	out.DailyWorkingDuration = parseHoursDuration(uuc.DailyWorkingDuration)
	out.AnnualVacationDays = uuc.AnnualVacationDays
	out.InitOvertimeDuration = parseHoursDuration(uuc.InitOvertimeDuration)
	out.InitVacationDays = uuc.InitVacationDays
	out.FirstWorkDay = parseDate(uuc.FirstWorkDay)
	return &out
}

// --- Role functions ---

// ToRoles converts a list of logic role models to an API roles model.
func ToRoles(roles []m.Role) *am.UserRoles {
	if roles == nil {
		return nil
	}

	rs := make([]string, len(roles))
	for i, role := range roles {
		rs[i] = role.String()
	}

	return &am.UserRoles{rs}
}

// FromRoles converts an API update roles model to a list of logic role models.
func FromRoles(uurs *am.UpdateUserRoles) []m.Role {
	if uurs == nil {
		return nil
	}

	rs := uurs.Roles
	roles := make([]m.Role, len(rs))
	for i, r := range rs {
		roles[i] = m.Role(r)
	}
	return roles
}
