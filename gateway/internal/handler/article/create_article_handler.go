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

func CreateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateArticleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article.NewCreateArticleLogic(r.Context(), svcCtx)
		err := l.CreateArticle(&req)
		response.Response(w, nil, err)
	}
}
