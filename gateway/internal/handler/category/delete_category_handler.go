// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package category

import (
	"net/http"

	"gateway/internal/logic/category"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := category.NewDeleteCategoryLogic(r.Context(), svcCtx)
		err := l.DeleteCategory(&req)
		response.Response(w, nil, err)
	}
}
