package article

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticlesByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByCategoryLogic {
	return &GetArticlesByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticlesByCategoryLogic) GetArticlesByCategory(req *types.GetArticlesByCategoryReq) (resp *types.ListArticlesData, err error) {
	r, err := l.svcCtx.Core.GetArticlesByCategory(l.ctx, &core_client.GetArticlesByCategoryReq{
		CategoryId: req.CategoryId,
		Page:       req.Page,
		PageSize:   req.PageSize,
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
