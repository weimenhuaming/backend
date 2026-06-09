// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"net/http"

	"gateway/internal/logic/article"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetArticleDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetArticleDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article.NewGetArticleDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetArticleDetail(&req)
		response.Response(w, resp, err)
	}
}
