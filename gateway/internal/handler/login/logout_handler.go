package login

import (
	"gateway/internal/utils"
	"net/http"

	"gateway/internal/logic/login"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewLogoutLogic(r.Context(), svcCtx)
		RefreshToken, err2 := utils.GetRefreshTokenFromRequest(r)
		if err2 != nil {
			return
		}
		resp, err := l.Logout(&req, RefreshToken)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
