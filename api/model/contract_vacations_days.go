package model

// ContractVacationDays
//
// Contains information about the monthly vacation days of a work contract.
//
// swagger:model ContractVacationDays
type ContractVacationDays struct {
	// The first day.
	// example: 2019-01-01
	FirstDay string `json:"firstDay"`

	// The number of days.
	// example: 26.0
	Days float32 `json:"days"`
}
