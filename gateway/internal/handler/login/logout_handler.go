// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"net/http"

	"gateway/internal/logic/login"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		refreshToken, err := utils.GetRefreshTokenFromRequest(r)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		l := login.NewLogoutLogic(r.Context(), svcCtx)
		err = l.Logout(&req, refreshToken)
		response.Response(w, nil, err)
	}
}
