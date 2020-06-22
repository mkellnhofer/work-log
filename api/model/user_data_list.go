package model

// UserList
//
// A list of users.
//
// swagger:model UserList
type UserDataList struct {
	// The users.
	Items []*UserData `json:"items"`
}

// NewUserDataList creates a new user list.
func NewUserDataList(items []*UserData) *UserDataList {
	return &UserDataList{items}
}
