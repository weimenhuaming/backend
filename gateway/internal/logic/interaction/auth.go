package interaction

import (
	"context"
	"strings"

	"gateway/internal/response"
	"gateway/internal/utils/vaild"
)

func userIDFromCtx(ctx context.Context) uint64 {
	userID, ok := vaild.GetUserID(ctx)
	if !ok {
		return 0
	}
	return userID
}

func likeUserFromCtx(ctx context.Context) (userID uint64, ok bool, msg string) {
	role, _ := ctx.Value(vaild.CtxUserRoleKey).(string)
	if role == "guest" {
		return 0, false, "游客无法操作"
	}
	userID, ok = vaild.GetUserID(ctx)
	if !ok {
		return 0, false, "请先登录"
	}
	return userID, true, ""
}

func authFailError(msg string) error {
	if msg == "游客无法操作" {
		return response.ErrorForbidden(msg)
	}
	return response.ErrorUnauthorized(msg)
}

func likeRPCError(err error) error {
	if err == nil {
		return response.ErrorInternalServer("unknown error")
	}
	msg := err.Error()
	if strings.Contains(msg, "已经点过赞") || strings.Contains(msg, "尚未点赞") {
		return response.ErrorBadRequest(msg)
	}
	return response.ErrorInternalServer(msg)
}
