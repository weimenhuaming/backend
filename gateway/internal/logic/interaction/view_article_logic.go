package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ViewArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewViewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewArticleLogic {
	return &ViewArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ViewArticleLogic) ViewArticle(req *types.ViewArticleReq) (resp *types.ViewArticleData, err error) {
	if req.ArticleId == 0 {
		return nil, response.ErrorBadRequest("文章ID无效")
	}

	r, err := l.svcCtx.Core.ViewArticle(l.ctx, &core_client.ViewArticleReq{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.ViewArticleData{ViewCount: r.GetViewCount()}, nil
}
