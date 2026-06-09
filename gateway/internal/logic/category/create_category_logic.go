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

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) error {
	if !vaild.IsAdmin(l.ctx) {
		return response.ErrorForbidden("非管理员，没有权限执行")
	}

	if req == nil || req.Name == "" {
		return response.ErrorBadRequest("name is required")
	}

	_, err := l.svcCtx.Core.CreateCategory(l.ctx, &core_client.CreateCategoryReq{
		Name: req.Name,
	})
	if err != nil {
		return response.ErrorInternalServer(err.Error())
	}

	return nil
}
