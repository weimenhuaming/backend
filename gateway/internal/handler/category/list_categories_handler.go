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

func ListCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCategoriesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := category.NewListCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.ListCategories(&req)
		response.Response(w, resp, err)
	}
}
