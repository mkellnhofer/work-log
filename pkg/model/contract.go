package model

import "time"

// ContractWorkingHours stores information about the daily working hours of a work contract.
type ContractWorkingHours struct {
	FirstDay time.Time // First day
	Hours    float32   // Number of hours
}

// ContractVacationDays stores information about the monthly vacation days of a work contract.
type ContractVacationDays struct {
	FirstDay time.Time // First day
	Days     float32   // Number of days
}

// Contract stores information about the work contract of a user.
type Contract struct {
	FirstDay          time.Time              // First day
	InitOvertimeHours float32                // Initial overtime hours
	InitVacationDays  float32                // Initial vacation days
	WorkingHours      []ContractWorkingHours // Daily working hours
	VacationDays      []ContractVacationDays // Monthly vacation days
}

// NewContract creates a new Contract model.
func NewContract() *Contract {
	return &Contract{}
}
