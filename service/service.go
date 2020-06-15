package service

import (
	"context"
	"fmt"

	"kellnhofer.com/work-log/constant"
	"kellnhofer.com/work-log/db/tx"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/model"
)

type service struct {
	tm *tx.TransactionManager
}

func getCurrentUserId(ctx context.Context) int {
	sc := ctx.Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	return sc.UserId
}

func checkHasCurrentUserRight(ctx context.Context, right model.Right) *e.Error {
	sc := ctx.Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	if hasUserRight(sc, right) {
		return nil
	}

	ec := getPermissionErrorCode(right)
	err := e.NewError(ec, fmt.Sprintf("User %d does not have required permission '%s'.", sc.UserId,
		right))
	log.Debug(err.StackTrace())
	return err
}

func hasCurrentUserRight(ctx context.Context, right model.Right) bool {
	sc := ctx.Value(constant.ContextKeySecurityContext).(*model.SecurityContext)
	return hasUserRight(sc, right)
}

func hasUserRight(sc *model.SecurityContext, right model.Right) bool {
	if sc.UserId == model.SystemUserId {
		return true
	}

	for _, ur := range sc.UserRoles {
		if hasRoleRight(ur, right) {
			return true
		}
	}

	return false
}

func hasRoleRight(role model.Role, right model.Right) bool {
	rrs := model.GetRoleRights(role)
	for _, rr := range rrs {
		if rr == right {
			return true
		}
	}
	return false
}

var permissionErrorCodes = map[model.Right]int{
	model.RightGetUserData:         e.PermGetUserData,
	model.RightChangeUserData:      e.PermChangeUserData,
	model.RightGetUserAccount:      e.PermGetUserAccount,
	model.RightChangeUserAccount:   e.PermChangeUserAccount,
	model.RightGetEntryCharacts:    e.PermGetEntryCharacts,
	model.RightChangeEntryCharacts: e.PermChangeEntryCharacts,
	model.RightGetAllEntries:       e.PermGetAllEntries,
	model.RightChangeAllEntries:    e.PermChangeAllEntries,
	model.RightGetOwnEntries:       e.PermGetOwnEntries,
	model.RightChangeOwnEntries:    e.PermChangeOwnEntries,
}

func getPermissionErrorCode(right model.Right) int {
	if permissionErrorCode, exists := permissionErrorCodes[right]; exists {
		return permissionErrorCode
	}
	return e.PermUnknown
}
