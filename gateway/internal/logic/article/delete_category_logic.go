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

type DeleteCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(req *types.DeleteCategoryReq) (resp *types.DeleteCategoryResp, err error) {
	_ = req
	if req == nil || req.Id == 0 {
		return &types.DeleteCategoryResp{
			Code: 400,
			Msg:  "id is required",
		}, nil
	}

	_, err = l.svcCtx.Core.DeleteCategory(l.ctx, &core_client.DeleteCategoryReq{
		Id: req.Id,
	})
	if err != nil {
		return &types.DeleteCategoryResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.DeleteCategoryResp{
		Code: 200,
		Msg:  "ok",
	}, nil
}
