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

func GetArticlesByCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetArticlesByCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article.NewGetArticlesByCategoryLogic(r.Context(), svcCtx)
		resp, err := l.GetArticlesByCategory(&req)
		response.Response(w, resp, err)
	}
}
