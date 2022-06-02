package model

// Contract
//
// Contains information about the work contract of a user.
//
// swagger:model Contract
type Contract struct {
	// The first work day.
	// example: 2019-01-01
	FirstDay string `json:"firstDay"`

	// The initial overtime hours.
	// example: 4.5
	InitOvertimeHours float32 `json:"initOvertimeHours"`

	// The initial vacation days.
	// example: 2.5
	InitVacationDays float32 `json:"initVacationDays"`

	// The daily working hours.
	WorkingHours []*ContractWorkingHours `json:"workingHours"`

	// The monthly vacation days.
	VacationDays []*ContractVacationDays `json:"vacationDays"`
}
