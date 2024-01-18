package model

// SecurityContext stores information about user who interacts with the application.
type SecurityContext struct {
	UserId    int    // ID of the current user
	UserRoles []Role // Roles of the current user
}

// NewSecurityContext creates a new SecurityContext model.
func NewSecurityContext(userId int, userRoles []Role) *SecurityContext {
	return &SecurityContext{userId, userRoles}
}

// IsSystemUser returns true for if this is the context for the system user.
func (sc SecurityContext) IsSystemUser() bool {
	return sc.UserId == SystemUserId
}

// IsAnonymousUser returns true for if this is a context for a anonymous user.
func (sc SecurityContext) IsAnonymousUser() bool {
	return sc.UserId == AnonymousUserId
}

// --- Helper functions ---

// GetSystemUserSecurityContext returns the security context for the system user.
func GetSystemUserSecurityContext() *SecurityContext {
	return NewSecurityContext(SystemUserId, []Role{RoleAdmin, RoleEvaluator, RoleUser})
}

// GetAnonymousUserSecurityContext returns the security context for a anonymous user.
func GetAnonymousUserSecurityContext() *SecurityContext {
	return NewSecurityContext(AnonymousUserId, []Role{})
}
