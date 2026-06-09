package category

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/response"
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

func (l *DeleteCategoryLogic) DeleteCategory(req *types.DeleteCategoryReq) error {
	_ = req

	role := l.ctx.Value("X-user-Role")
	if role != "admin" {
		return response.NewError(403, "非管理员，没有权限执行")
	}

	if req == nil || req.Id == 0 {
		return response.NewError(400, "id is required")
	}

	_, err := l.svcCtx.Core.DeleteCategory(l.ctx, &core_client.DeleteCategoryReq{
		Id: req.Id,
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
