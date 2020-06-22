package model

// UpdateUser
//
// Holds the new information about a user.
//
// swagger:model UpdateUser
type UpdateUserData struct {
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
	Contract *UpdateUserContract `json:"contract"`
}
