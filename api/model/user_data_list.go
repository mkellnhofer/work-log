package model

// UserDataList holds a list of users.
type UserDataList struct {
	Items []*UserData `json:"items"`
}

// NewUserList creates a new user list.
func NewUserList(items []*UserData) *UserDataList {
	return &UserDataList{items}
}
