package logic

import (
	"context"
	"core-rpc/internal/utils"
	"strings"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchArticlesLogic) SearchArticles(in *core.SearchArticlesReq) (*core.SearchArticlesResp, error) {
	keyword := strings.TrimSpace(in.Keyword)
	page := utils.NormalizePageUint32(in.Page)
	size := utils.NormalizePageSizeUint32(in.PageSize, 10)
	off, limit := utils.OffsetLimit(page, size)

	articles, total, err := l.svcCtx.ArtRepo.Search(keyword, in.CategoryId, off, limit)
	if err != nil {
		return nil, err
	}

	protoList, err := l.svcCtx.ArtRepo.LoadArticlesWithAuthors(articles)
	if err != nil {
		return nil, err
	}

	return &core.SearchArticlesResp{
		Articles: protoList,
		Total:    uint32(total),
		Page:     uint32(page),
		PageSize: uint32(size),
	}, nil
}
