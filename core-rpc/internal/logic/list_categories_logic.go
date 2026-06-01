package logic

import (
	"context"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCategoriesLogic) ListCategories(in *core.ListCategoriesReq) (*core.ListCategoriesResp, error) {
	var categories []entity.Category
	if err := l.svcCtx.Db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	list := make([]*core.CategoryInfo, 0, len(categories))
	for _, c := range categories {
		list = append(list, &core.CategoryInfo{
			Id:   c.ID,
			Name: c.Name,
		})
	}
	return &core.ListCategoriesResp{Categories: list}, nil
}
