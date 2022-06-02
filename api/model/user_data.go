package model

// User
//
// Contains information about a user.
//
// swagger:model User
type UserData struct {
	// The ID of the user.
	// example: 1
	Id int `json:"id"`

	// The name of the user.
	// min length: 1
	// max length: 100
	// example: John
	Name string `json:"name"`

	// The username of the user.
	// min length: 1
	// max length: 100
	// example: john
	Username string `json:"username"`

	// The work contract of the user.
	Contract *Contract `json:"contract"`
}
