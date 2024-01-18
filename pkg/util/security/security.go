package security

import (
	"context"

	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/model"
)

// CreateSystemContext returns a context with system user information
func CreateSystemContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.ContextKeySecurityContext,
		model.GetSystemUserSecurityContext())
}

// GetSecurityContext returns the security context.
func GetSecurityContext(ctx context.Context) *model.SecurityContext {
	return ctx.Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
}

// GetCurrentUserId returns the user ID from the security context.
func GetCurrentUserId(ctx context.Context) int {
	sc := GetSecurityContext(ctx)
	return sc.UserId
}
