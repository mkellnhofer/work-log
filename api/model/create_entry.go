package model

// CreateEntry contains information about a entry.
type CreateEntry struct {
	UserId        int    `json:"userId"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	BreakDuration int    `json:"breakDuration"`
	TypeId        int64  `json:"typeId"`
	ActivityId    int64  `json:"activityId"`
	Desciption    string `json:"description"`
}
