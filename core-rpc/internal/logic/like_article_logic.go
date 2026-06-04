package logic

import (
	"context"
	"errors"

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

	// verify article exists
	if _, err := l.svcCtx.ArtRepo.FindByID(in.ArticleId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	likeCount, err := l.svcCtx.InteractionRepo.LikeArticle(in.UserId, in.ArticleId)
	if err != nil {
		return nil, err
	}
	return &core.LikeArticleResp{LikeCount: likeCount}, nil
}
