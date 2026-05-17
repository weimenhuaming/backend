package article

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteArticleLogic {
	return &FavoriteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteArticleLogic) FavoriteArticle(req *types.FavorArticleReq) (resp *types.FavorArticleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
