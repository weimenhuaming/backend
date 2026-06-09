package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLikedArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLikedArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLikedArticlesLogic {
	return &ListUserLikedArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLikedArticlesLogic) ListUserLikedArticles(req *types.ListUserLikedArticlesReq) (resp *types.ListUserLikedArticlesData, err error) {
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

	r, err := l.svcCtx.Core.ListUserLikedArticles(l.ctx, &core_client.ListUserLikedArticlesReq{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	return &types.ListUserLikedArticlesData{
		Articles: toTypesArticleList(r.GetArticles()),
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
