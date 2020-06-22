package model

// CreateEntry
//
// Holds information about a new entry.
//
// swagger:model CreateEntry
type CreateEntry struct {
	// The ID of the user.
	// example: 1
	UserId int `json:"userId"`

	// The start time of the entry.
	// example: 2019-01-01T15:00:00
	StartTime string `json:"startTime"`

	// The end time of the entry.
	// example: 2019-01-01T16:00:00
	EndTime string `json:"endTime"`

	// The break duration of the entry in minutes.
	// example: 60
	BreakDuration int `json:"breakDuration"`

	// The ID of the entry type.
	// example: 1
	TypeId int `json:"typeId"`

	// The ID of the entry activity.
	// example: 1
	ActivityId int `json:"activityId"`

	// The description with additional information about the entry.
	// min length: 0
	// max length: 200
	Description string `json:"description"`
}
