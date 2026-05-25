package logic

import (
	"context"

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
	cats, err := l.svcCtx.CategoryModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}

	resp := &core.ListCategoriesResp{}
	for _, c := range cats {
		resp.Categories = append(resp.Categories, &core.CategoryInfo{
			Id:   c.Id,
			Name: c.Name,
		})
	}

	return resp, nil
}
