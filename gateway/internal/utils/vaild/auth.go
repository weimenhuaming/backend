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

func GetActionUserID(ctx context.Context) (userID uint64, ok bool, msg string) {
	role, _ := ctx.Value(CtxUserRoleKey).(string)
	if role == "guest" {
		return 0, false, "游客无法操作"
	}
	userID, ok = GetUserID(ctx)
	if !ok {
		return 0, false, "请先登录"
	}
	return userID, true, ""
}
