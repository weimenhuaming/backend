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

type ViewArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewViewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ViewArticleLogic {
	return &ViewArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ViewArticleLogic) ViewArticle(in *core.ViewArticleReq) (*core.ViewArticleResp, error) {
	if in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}

	var article entity.Article
	if err := l.svcCtx.Db.First(&article, in.ArticleId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	if err := l.svcCtx.Db.Model(&entity.Article{}).Where("id = ?", in.ArticleId).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		return nil, err
	}

	if err := l.svcCtx.Db.Select("view_count").First(&article, in.ArticleId).Error; err != nil {
		return nil, err
	}

	return &core.ViewArticleResp{ViewCount: article.ViewCount}, nil
}
