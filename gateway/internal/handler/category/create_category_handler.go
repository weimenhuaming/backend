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

func CreateCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := category.NewCreateCategoryLogic(r.Context(), svcCtx)
		err := l.CreateCategory(&req)
		response.Response(w, nil, err)
	}
}
