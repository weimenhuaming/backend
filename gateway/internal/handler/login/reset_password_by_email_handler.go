// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"net/http"

	"gateway/internal/logic/login"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Reset_password_by_emailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewReset_password_by_emailLogic(r.Context(), svcCtx)
		err := l.Reset_password_by_email(&req)
		response.Response(w, nil, err)
	}
}
