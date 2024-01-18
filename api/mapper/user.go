package mapper

import (
	"strings"

	am "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/pkg/model"
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
	out.Contract = toContract(ud.Contract)
	return &out
}

// FromCreateUserData converts an API user creation model to a logic user data model.
func FromCreateUserData(cud *am.CreateUserData) *m.UserData {
	if cud == nil {
		return nil
	}

	var out m.UserData
	out.User = fromCreateUser(cud)
	out.Contract = fromCreateContract(cud.Contract)
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
	out.Contract = fromUpdateContract(uud.Contract)
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

func toContract(uc *m.Contract) *am.Contract {
	if uc == nil {
		return nil
	}

	var out am.Contract
	out.FirstDay = formatDate(uc.FirstDay)
	out.InitOvertimeHours = uc.InitOvertimeHours
	out.InitVacationDays = uc.InitVacationDays
	out.WorkingHours = toContractWorkingHours(uc.WorkingHours)
	out.VacationDays = toContractVacationDays(uc.VacationDays)
	return &out
}

func fromCreateContract(cuc *am.CreateContract) *m.Contract {
	if cuc == nil {
		return nil
	}

	var out m.Contract
	out.FirstDay = parseDate(cuc.FirstDay)
	out.InitOvertimeHours = cuc.InitOvertimeHours
	out.InitVacationDays = cuc.InitVacationDays
	out.WorkingHours = fromContractWorkingHours(cuc.WorkingHours)
	out.VacationDays = fromContractVacationDays(cuc.VacationDays)
	return &out
}

func fromUpdateContract(uuc *am.UpdateContract) *m.Contract {
	if uuc == nil {
		return nil
	}

	var out m.Contract
	out.FirstDay = parseDate(uuc.FirstDay)
	out.InitOvertimeHours = uuc.InitOvertimeHours
	out.InitVacationDays = uuc.InitVacationDays
	out.WorkingHours = fromContractWorkingHours(uuc.WorkingHours)
	out.VacationDays = fromContractVacationDays(uuc.VacationDays)
	return &out
}

func toContractWorkingHours(whs []m.ContractWorkingHours) []*am.ContractWorkingHours {
	outs := make([]*am.ContractWorkingHours, len(whs))
	for i, wh := range whs {
		outs[i] = &am.ContractWorkingHours{}
		outs[i].FirstDay = formatDate(wh.FirstDay)
		outs[i].Hours = wh.Hours
	}
	return outs
}

func fromContractWorkingHours(whs []*am.ContractWorkingHours) []m.ContractWorkingHours {
	outs := make([]m.ContractWorkingHours, len(whs))
	for i, wh := range whs {
		if wh != nil {
			outs[i].FirstDay = parseDate(wh.FirstDay)
			outs[i].Hours = wh.Hours
		}
	}
	return outs
}

func toContractVacationDays(vds []m.ContractVacationDays) []*am.ContractVacationDays {
	outs := make([]*am.ContractVacationDays, len(vds))
	for i, vd := range vds {
		outs[i] = &am.ContractVacationDays{}
		outs[i].FirstDay = formatDate(vd.FirstDay)
		outs[i].Days = vd.Days
	}
	return outs
}

func fromContractVacationDays(vds []*am.ContractVacationDays) []m.ContractVacationDays {
	outs := make([]m.ContractVacationDays, len(vds))
	for i, vd := range vds {
		if vd != nil {
			outs[i].FirstDay = parseDate(vd.FirstDay)
			outs[i].Days = vd.Days
		}
	}
	return outs
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
