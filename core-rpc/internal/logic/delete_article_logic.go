package logic

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleLogic) DeleteArticle(in *core.DeleteArticleReq) (*core.DeleteArticleResp, error) {
	uid := in.GetUserId()
	if uid == 0 {
		return nil, errors.New("missing user id")
	}

	a, err := l.svcCtx.ArticleModel.FindOne(l.ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	if a.UserId != uid {
		return nil, errors.New("not article owner")
	}

	err = l.svcCtx.ArticleModel.SoftDelete(l.ctx, in.GetId(), uid)
	if err != nil {
		return nil, err
	}

	return &core.DeleteArticleResp{}, nil
}
