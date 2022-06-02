package model

// ContractWorkingHours
//
// Contains information about the daily working hours of a work contract.
//
// swagger:model ContractWorkingHours
type ContractWorkingHours struct {
	// The first day.
	// example: 2019-01-01
	FirstDay string `json:"firstDay"`

	// The number of hours.
	// example: 8.0
	Hours float32 `json:"hours"`
}
