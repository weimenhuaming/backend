// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"net/http"

	"gateway/internal/logic/upload"
	"gateway/internal/response"
	"gateway/internal/svc"
)

func UploadAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := upload.NewUploadAvatarLogic(r.Context(), svcCtx)
		resp, err := l.UploadAvatar(r)
		response.Response(w, resp, err)
	}
}
