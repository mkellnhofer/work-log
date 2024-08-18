package model

// UpdateEntry
//
// Holds the new information about a entry.
//
// swagger:model UpdateEntry
type UpdateEntry struct {
	// The ID of the user.
	// example: 1
	UserId int `json:"userId"`

	// The start time of the entry.
	// example: 2019-01-01T15:00:00
	StartTime string `json:"startTime"`

	// The end time of the entry.
	// example: 2019-01-01T16:00:00
	EndTime string `json:"endTime"`

	// The ID of the entry type.
	// example: 1
	TypeId int `json:"typeId"`

	// The ID of the entry activity.
	// example: 1
	ActivityId int `json:"activityId"`

	// The name of the project.
	// min length: 0
	// max length: 30
	// example: Web Client
	Project string `json:"project"`

	// The description with additional information about the entry.
	// min length: 0
	// max length: 200
	Description string `json:"description"`

	// The labels associated with the entry.
	// min length: 3
	// max length: 20
	// example: ["bug", "frontend"]
	Labels []string `json:"labels"`
}
