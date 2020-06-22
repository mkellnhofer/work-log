package model

// UserRoles
//
// A list of roles of a user.
//
// swagger:model UserRoles
type UserRoles struct {
	// The roles.
	Roles []string `json:"roles"`
}
