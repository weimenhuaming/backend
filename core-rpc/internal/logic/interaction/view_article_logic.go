package interaction

import (
	"context"
	"errors"

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

	a, err := l.svcCtx.ArtRepo.FindByID(in.ArticleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	if err := l.svcCtx.ArtRepo.IncView(in.ArticleId); err != nil {
		return nil, err
	}

	// fetch updated count
	a, err = l.svcCtx.ArtRepo.FindByID(in.ArticleId)
	if err != nil {
		return nil, err
	}

	return &core.ViewArticleResp{ViewCount: a.ViewCount}, nil
}
