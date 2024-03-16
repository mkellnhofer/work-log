package service

import (
	"context"

	"kellnhofer.com/work-log/pkg/db/tx"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/util/security"
)

type service struct {
	tm *tx.TransactionManager
}

func getCurrentUserId(ctx context.Context) int {
	return security.GetCurrentUserId(ctx)
}

func checkHasCurrentUserRight(ctx context.Context, right model.Right) error {
	return security.CheckHasCurrentUserRight(ctx, right)
}

func hasCurrentUserRight(ctx context.Context, right model.Right) bool {
	return security.HasCurrentUserRight(ctx, right)
}
