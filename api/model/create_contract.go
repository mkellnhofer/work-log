package model

// CreateContract
//
// Holds information about the work contract of a new user.
//
// swagger:model CreateContract
type CreateContract struct {
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
