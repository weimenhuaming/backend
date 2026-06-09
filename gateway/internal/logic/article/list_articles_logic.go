package article

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

type ListArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListArticlesLogic) ListArticles(req *types.ListArticlesReq) (resp *types.ListArticlesData, err error) {
	page, pageSize := vaild.NormalizePageSize(req.Page, req.PageSize)

	r, err := l.svcCtx.Core.ListArticles(l.ctx, &core_client.ListArticlesReq{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.ListArticlesData{
		Articles: converter.ToArticleList(r.GetArticles()),
		Total:    r.GetTotal(),
		Page:     r.GetPage(),
		PageSize: r.GetPageSize(),
	}, nil
}
