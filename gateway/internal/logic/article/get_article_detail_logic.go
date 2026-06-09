package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(req *types.GetArticleDetailReq) (resp *types.ArticleInfo, err error) {
	// call core rpc to get article detail
	r, err := l.svcCtx.Core.GetArticleDetail(l.ctx, &core_client.GetArticleDetailReq{Id: req.Id})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	if r == nil || r.Article == nil {
		return nil, response.ErrorNotFound("article not found")
	}

	info := converter.ToArticleInfo(r.GetArticle())
	return &info, nil
}
