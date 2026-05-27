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

type LikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeArticleLogic) LikeArticle(in *core.LikeArticleReq) (*core.LikeArticleResp, error) {
	if in.UserId == 0 || in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}

	var count int64
	if err := l.svcCtx.Db.Model(&entity.Article{}).Where("id = ?", in.ArticleId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("文章不存在")
	}

	var likeCount uint32
	err := l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		delta, err := toggleArticleLike(tx, in.UserId, in.ArticleId, in.IsLike)
		if err != nil {
			return err
		}
		likeCount, err = adjustArticleCounter(tx, in.ArticleId, "like_count", delta)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &core.LikeArticleResp{LikeCount: likeCount}, nil
}
