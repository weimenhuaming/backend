// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package agent

import (
	"net/http"

	"gateway/internal/logic/agent"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSessionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewCreateSessionLogic(r.Context(), svcCtx)
		resp, err := l.CreateSession(&req)
		response.Response(w, resp, err)
	}
}
