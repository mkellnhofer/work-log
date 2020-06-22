package model

// UpdateUserRoles
//
// Holds the new list of roles of a user.
//
// swagger:model UpdateUserRoles
type UpdateUserRoles struct {
	// The roles.
	Roles []string `json:"roles"`
}
