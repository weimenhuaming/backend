package login

import (
	"net/http"
	"time"

	"gateway/internal/logic/login"
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
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// 设置加上cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    resp.Data.RefreshToken,
				HttpOnly: true,
				Path:     "/", // 所有路径下都生效
				Expires:  time.Now().Add(time.Duration(svcCtx.Config.RefreshExpire) * time.Second),
				SameSite: http.SameSiteLaxMode,
				Secure:   false, // 如果你是 HTTPS，必须设置
				// http本地测试需要关掉
			})
			// 不应该把refresh token传递给前端，应该在后端设置cookie
			resp.Data.RefreshToken = ""
			w.Header().Set("Authorization", resp.Data.AccessToken)
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
