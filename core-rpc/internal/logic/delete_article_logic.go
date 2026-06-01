package logic

import (
	"context"
	"errors"

	"core-rpc/internal/model/entity"
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
	var article entity.Article
	if err := l.svcCtx.Db.First(&article, in.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}
	if in.UserId > 0 && article.UserID != in.UserId {
		return nil, errors.New("无权删除该文章")
	}
	if err := l.svcCtx.Db.Delete(&article).Error; err != nil {
		return nil, err
	}
	return &core.DeleteArticleResp{}, nil
}
