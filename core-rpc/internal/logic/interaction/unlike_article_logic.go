package interaction

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UnlikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeArticleLogic {
	return &UnlikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnlikeArticleLogic) UnlikeArticle(in *core.UnlikeArticleReq) (*core.UnlikeArticleResp, error) {
	if in.UserId == 0 || in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}

	if _, err := l.svcCtx.ArtRepo.FindByID(in.ArticleId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	likeCount, err := l.svcCtx.InteractionRepo.UnlikeArticle(in.UserId, in.ArticleId)
	if err != nil {
		return nil, err
	}
	return &core.UnlikeArticleResp{LikeCount: likeCount}, nil
}
