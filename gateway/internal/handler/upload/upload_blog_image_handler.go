package upload

import (
	"net/http"

	"gateway/internal/logic/upload"
	"gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadBlogImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := upload.NewUploadBlogImageLogic(r.Context(), svcCtx)
		resp, err := l.UploadBlogImage(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
