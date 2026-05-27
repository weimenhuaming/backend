package logic

import (
	"context"
	"strings"

	"core-rpc/internal/model/entity"
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
	page := normalizePageUint32(in.Page)
	size := normalizePageSizeUint32(in.PageSize, 10)
	off, limit := offsetLimit(page, size)

	q := l.svcCtx.Db.Model(&entity.Article{})
	if in.CategoryId > 0 {
		q = q.Where("category_id = ?", in.CategoryId)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	q = q.Order("created_at DESC")

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

	return &core.SearchArticlesResp{
		Articles: protoList,
		Total:    uint32(total),
		Page:     uint32(page),
		PageSize: uint32(size),
	}, nil
}
