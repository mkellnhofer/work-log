package model

// UserInfo stores basic view data for a user.
type UserInfo struct {
	Id       int
	Initials string
}

// UserProfileInfo stores detailed view data for a user.
type UserProfileInfo struct {
	Id       int
	Initials string
	Name string
	Username string
	Contract *ContractInfo
}
