package logic

import (
	"context"
	"core-rpc/internal/utils"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailLogic {
	return &GetArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleDetailLogic) GetArticleDetail(in *core.GetArticleDetailReq) (*core.GetArticleDetailResp, error) {
	var article entity.Article
	if err := l.svcCtx.Db.First(&article, in.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	_ = l.svcCtx.Db.Model(&article).Update("view_count", gorm.Expr("view_count + ?", 1))
	article.ViewCount++

	var author entity.User
	authorName, authorAvatar := "", ""
	if err := l.svcCtx.Db.Select("name", "avatar").First(&author, article.UserID).Error; err == nil {
		authorName = author.Name
		authorAvatar = author.Avatar
	}

	return &core.GetArticleDetailResp{
		Article: utils.ArticleToProto(&article, authorName, authorAvatar),
	}, nil
}
