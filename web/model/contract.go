package model

// ContractInfo stores view data of the user contract.
type ContractInfo struct {
	FirstDay          string
	InitOvertimeHours string
	InitVacationDays  string
	WorkingHours      []*ContractWorkingHours
	VacationDays      []*ContractVacationDays
}

// ContractWorkingHours stores view data of the user contract working hours.
type ContractWorkingHours struct {
	FirstDay string
	Hours    string
}

// ContractVacationDays stores view data of the user contract vacation days.
type ContractVacationDays struct {
	FirstDay string
	Days     string
}
