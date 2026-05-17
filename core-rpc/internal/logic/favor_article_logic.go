package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavorArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFavorArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavorArticleLogic {
	return &FavorArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FavorArticleLogic) FavorArticle(in *core.FavoriteArticleReq) (*core.FavoriteArticleResp, error) {
	// todo: add your logic here and delete this line

	return &core.FavoriteArticleResp{}, nil
}
