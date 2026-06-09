package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
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

func (l *ListMyArticlesLogic) ListMyArticles(req *types.ListMyArticlesReq) (resp *types.ListMyArticlesData, err error) {
	if currentUserRole(l.ctx) != "admin" {
		return nil, response.NewError(403, "非管理员，没有权限执行")
	}

	userId, ok := currentUserID(l.ctx)
	if !ok {
		return nil, response.NewError(401, "用户未登录")
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
		return nil, response.NewError(500, err.Error())
	}

	return &types.ListMyArticlesData{
		Articles: toTypesArticleList(r.GetArticles()),
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
