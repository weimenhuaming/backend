package category

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/response"
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

func (l *ListCategoriesLogic) ListCategories(req *types.ListCategoriesReq) (resp *types.ListCategoriesData, err error) {
	_ = req
	rpcResp, err := l.svcCtx.Core.ListCategories(l.ctx, &core_client.ListCategoriesReq{})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	var cats []types.CategoryInfo
	for _, c := range rpcResp.GetCategories() {
		cats = append(cats, types.CategoryInfo{
			Id:   c.GetId(),
			Name: c.GetName(),
		})
	}

	return &types.ListCategoriesData{Categories: cats}, nil
}
