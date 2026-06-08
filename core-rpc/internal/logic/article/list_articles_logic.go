package article

import (
	"context"
	"core-rpc/internal/utils"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesLogic {
	return &ListArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticlesLogic) ListArticles(in *core.ListArticlesReq) (*core.ListArticlesResp, error) {
	page := utils.NormalizePageUint32(in.Page)
	size := utils.NormalizePageSizeUint32(in.PageSize, 10)
	off, limit := utils.OffsetLimit(page, size)

	articles, total, err := l.svcCtx.ArtRepo.List(off, limit)
	if err != nil {
		return nil, err
	}

	protoList, err := l.svcCtx.ArtRepo.LoadArticlesWithAuthors(articles)
	if err != nil {
		return nil, err
	}

	return &core.ListArticlesResp{
		Articles: protoList,
		Total:    uint32(total),
		Page:     uint32(page),
		PageSize: uint32(size),
	}, nil
}
