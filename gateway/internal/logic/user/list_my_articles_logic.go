package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"
	"gateway/internal/utils/vaild"

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
	if !vaild.IsAdmin(l.ctx) {
		return nil, response.ErrorForbidden("非管理员，没有权限执行")
	}

	userId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return nil, response.ErrorUnauthorized("用户未登录")
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
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.ListMyArticlesData{
		Articles: converter.ToArticleList(r.GetArticles()),
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
