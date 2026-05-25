// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategoriesLogic) ListCategories(req *types.ListCategoriesReq) (resp *types.ListCategoriesResp, err error) {
	_ = req
	rpcResp, err := l.svcCtx.Core.ListCategories(l.ctx, &core_client.ListCategoriesReq{})
	if err != nil {
		return &types.ListCategoriesResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	var cats []types.CategoryInfo
	for _, c := range rpcResp.GetCategories() {
		cats = append(cats, types.CategoryInfo{
			Id:   c.GetId(),
			Name: c.GetName(),
		})
	}

	return &types.ListCategoriesResp{
		Code: 200,
		Msg:  "ok",
		Data: types.ListCategoriesData{Categories: cats},
	}, nil
}
