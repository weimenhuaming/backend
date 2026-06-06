package login

import (
	"net/http"
	"time"

	"gateway/internal/logic/login"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.Data.RefreshToken,
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(time.Duration(svcCtx.Config.RefreshExpire) * time.Second),
				SameSite: http.SameSiteLaxMode,
				Secure:   false,
			})
			resp.Data.RefreshToken = ""
			w.Header().Set("Authorization", resp.Data.AccessToken)
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
