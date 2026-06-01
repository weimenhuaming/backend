package interaction

import (
	"context"
	"strings"
)

func userIDFromCtx(ctx context.Context) uint64 {
	userID, ok := ctx.Value("X-user-Id").(uint64)
	if !ok {
		return 0
	}
	return userID
}

func likeUserFromCtx(ctx context.Context) (userID uint64, ok bool, msg string) {
	role, _ := ctx.Value("X-user-Role").(string)
	if role == "guest" {
		return 0, false, "游客无法操作"
	}
	userID, ok = ctx.Value("X-user-Id").(uint64)
	if !ok || userID == 0 {
		return 0, false, "请先登录"
	}
	return userID, true, ""
}

func authFailCode(msg string) int {
	if msg == "游客无法操作" {
		return 403
	}
	return 401
}

func likeErrCode(err error) int {
	if err == nil {
		return 500
	}
	msg := err.Error()
	if strings.Contains(msg, "已经点过赞") || strings.Contains(msg, "尚未点赞") {
		return 400
	}
	return 500
}
