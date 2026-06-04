package logic

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	art, err := l.svcCtx.ArtRepo.FindByID(in.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	if in.UserId > 0 && art.UserID != in.UserId {
		return nil, errors.New("无权删除该文章")
	}
	if err := l.svcCtx.ArtRepo.DeleteByID(in.Id); err != nil {
		return nil, err
	}
	return &core.DeleteArticleResp{}, nil
}
