// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interaction

import (
	"net/http"

	"gateway/internal/logic/interaction"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LikeArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeArticleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := interaction.NewLikeArticleLogic(r.Context(), svcCtx)
		resp, err := l.LikeArticle(&req)
		response.Response(w, resp, err)
	}
}
