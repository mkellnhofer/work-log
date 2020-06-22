package model

// UpdateEntry contains information about a entry.
type UpdateEntry struct {
	UserId        int    `json:"userId"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	BreakDuration int    `json:"breakDuration"`
	TypeId        int64  `json:"typeId"`
	ActivityId    int64  `json:"activityId"`
	Desciption    string `json:"description"`
}
