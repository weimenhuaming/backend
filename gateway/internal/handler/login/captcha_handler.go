// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"net/http"

	"gateway/internal/logic/login"
	"gateway/internal/response"
	"gateway/internal/svc"
)

func CaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := login.NewCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.Captcha()
		response.Response(w, resp, err)
	}
}
