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

func SendEmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewSendEmailLogic(r.Context(), svcCtx)
		err := l.SendEmail(&req)
		response.Response(w, nil, err)
	}
}
