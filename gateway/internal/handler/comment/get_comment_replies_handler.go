// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"net/http"

	"gateway/internal/logic/comment"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCommentRepliesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCommentRepliesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := comment.NewGetCommentRepliesLogic(r.Context(), svcCtx)
		resp, err := l.GetCommentReplies(&req)
		response.Response(w, resp, err)
	}
}
