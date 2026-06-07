package agent

import "context"

func requireAdmin(ctx context.Context) (code int, msg string, ok bool) {
	role, _ := ctx.Value("X-user-Role").(string)
	if role != "admin" {
		return 403, "非管理员，没有权限执行", false
	}
	return 0, "", true
}
