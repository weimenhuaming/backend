package category

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (resp *types.CreateCategoryResp, err error) {
	if req == nil || req.Name == "" {
		return &types.CreateCategoryResp{
			Code: 400,
			Msg:  "name is required",
		}, nil
	}

	_, err = l.svcCtx.Core.CreateCategory(l.ctx, &core_client.CreateCategoryReq{
		Name: req.Name,
	})
	if err != nil {
		return &types.CreateCategoryResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.CreateCategoryResp{
		Code: 200,
		Msg:  "ok",
	}, nil

	return
}
