package vaild

import (
	"context"
)

const (
	CtxUserIDKey   = "X-user-Id"
	CtxUserRoleKey = "X-user-Role"
)

func IsAdmin(ctx context.Context) bool {
	role, _ := ctx.Value(CtxUserRoleKey).(string)
	if role != "admin" {
		return false
	}
	return true
}

func GetUserID(ctx context.Context) (uint64, bool) {
	userId, ok := ctx.Value(CtxUserIDKey).(uint64)
	return userId, ok && userId > 0
}
