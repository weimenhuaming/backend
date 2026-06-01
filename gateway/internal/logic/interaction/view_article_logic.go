package interaction

import (
	"context"

	core_client "core-rpc/core_client"
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

func (l *ViewArticleLogic) ViewArticle(req *types.ViewArticleReq) (resp *types.ViewArticleResp, err error) {
	if req.ArticleId == 0 {
		return &types.ViewArticleResp{Code: 400, Msg: "文章ID无效"}, nil
	}

	r, err := l.svcCtx.Core.ViewArticle(l.ctx, &core_client.ViewArticleReq{
		ArticleId: req.ArticleId,
	})
	if err != nil {
		return &types.ViewArticleResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.ViewArticleResp{
		Code: 200,
		Msg:  "ok",
		Data: types.ViewArticleData{ViewCount: r.GetViewCount()},
	}, nil
}
