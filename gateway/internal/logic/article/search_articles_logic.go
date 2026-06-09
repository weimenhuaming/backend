package article

import (
	"context"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticlesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchArticlesLogic) SearchArticles(req *types.SearchArticlesReq) (resp *types.SearchArticlesData, err error) {
	// todo: add your logic here and delete this line

	return
}
