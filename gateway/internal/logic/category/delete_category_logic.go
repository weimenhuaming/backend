package category

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/vaild"

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
	if _, ok, msg := vaild.GetAdminUserID(l.ctx); !ok {
		return response.ErrorAdminAuth(msg)
	}
	if req == nil || req.Id == 0 {
		return response.ErrorBadRequest("id is required")
	}

	_, err := l.svcCtx.Core.DeleteCategory(l.ctx, &core_client.DeleteCategoryReq{
		Id: req.Id,
	})
	if err != nil {
		return response.ErrorInternalServer(err.Error())
	}

	return nil
}
