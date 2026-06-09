// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"gateway/internal/logic/user"
	"gateway/internal/response"
	"gateway/internal/svc"
)

func GetUserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetUserProfile()
		response.Response(w, resp, err)
	}
}
