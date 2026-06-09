// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package agent

import (
	"net/http"

	"gateway/internal/logic/agent"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteKnowledgeCollectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteKnowledgeCollectionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := agent.NewDeleteKnowledgeCollectionLogic(r.Context(), svcCtx)
		err := l.DeleteKnowledgeCollection(&req)
		response.Response(w, nil, err)
	}
}
