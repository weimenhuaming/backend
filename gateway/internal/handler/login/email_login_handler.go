// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"net/http"
	"time"

	"gateway/internal/logic/login"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func EmailLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginEmailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewEmailLoginLogic(r.Context(), svcCtx)
		resp, err := l.EmailLogin(&req)
		if err != nil {
			response.Response(w, nil, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    resp.RefreshToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Duration(svcCtx.Config.RefreshExpire) * time.Second),
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		})
		resp.RefreshToken = ""
		w.Header().Set("Authorization", resp.AccessToken)
		response.Response(w, resp, nil)
	}
}
