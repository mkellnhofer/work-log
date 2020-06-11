package security

import (
	"context"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/model"
)

// CreateSystemContext returns a context with system user information
func CreateSystemContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.ContextKeySecurityContext,
		model.GetSystemUserSecurityContext())
}
