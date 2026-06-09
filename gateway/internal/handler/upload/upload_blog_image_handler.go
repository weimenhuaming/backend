// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"net/http"

	"gateway/internal/logic/upload"
	"gateway/internal/response"
	"gateway/internal/svc"
)

func UploadBlogImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := upload.NewUploadBlogImageLogic(r.Context(), svcCtx)
		resp, err := l.UploadBlogImage(r)
		response.Response(w, resp, err)
	}
}
