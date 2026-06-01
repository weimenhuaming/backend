package logic

import (
	"context"
	"core-rpc/internal/utils"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticlesByCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticlesByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticlesByCategoryLogic {
	return &GetArticlesByCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticlesByCategoryLogic) GetArticlesByCategory(in *core.GetArticlesByCategoryReq) (*core.ListArticlesResp, error) {
	page := utils.NormalizePageUint32(in.Page)
	size := utils.NormalizePageSizeUint32(in.PageSize, 10)
	off, limit := utils.OffsetLimit(page, size)

	q := l.svcCtx.Db.Model(&entity.Article{}).Where("category_id = ?", in.CategoryId).Order("created_at DESC")

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
