package logic

import (
	"context"

	"core-rpc/internal/model/entity"
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
	page := normalizePageUint32(in.Page)
	size := normalizePageSizeUint32(in.PageSize, 10)
	off, limit := offsetLimit(page, size)

	q := listArticlesQuery(l.svcCtx.Db, in.CategoryId, in.UserId, in.SortBy, in.SortOrder)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	var articles []entity.Article
	if err := q.Offset(off).Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}

	protoList, err := loadArticlesWithAuthors(l.svcCtx.Db, articles)
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
