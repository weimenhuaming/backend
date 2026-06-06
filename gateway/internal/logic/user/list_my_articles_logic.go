package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyArticlesLogic {
	return &ListMyArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyArticlesLogic) ListMyArticles(req *types.ListMyArticlesReq) (resp *types.ListMyArticlesResp, err error) {
	if currentUserRole(l.ctx) != "admin" {
		return &types.ListMyArticlesResp{Code: 403, Msg: "非管理员，没有权限执行"}, nil
	}

	userId, ok := currentUserID(l.ctx)
	if !ok {
		return &types.ListMyArticlesResp{Code: 401, Msg: "用户未登录"}, nil
	}

	page := req.Page
	pageSize := req.PageSize
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	r, err := l.svcCtx.Core.ListUserArticles(l.ctx, &core_client.ListUserArticlesReq{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return &types.ListMyArticlesResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.ListMyArticlesResp{
		Code: 200,
		Msg:  "ok",
		Data: types.ListMyArticlesData{
			Articles: toTypesArticleList(r.GetArticles()),
			Total:    r.GetTotal(),
			Page:     r.GetPage(),
			PageSize: r.GetPageSize(),
		},
	}, nil
}
