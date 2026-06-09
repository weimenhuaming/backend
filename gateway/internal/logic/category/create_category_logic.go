package category

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/response"
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

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) error {
	role := l.ctx.Value("X-user-Role")
	if role != "admin" {
		return response.NewError(403, "非管理员，没有权限执行")
	}

	if req == nil || req.Name == "" {
		return response.NewError(400, "name is required")
	}

	_, err := l.svcCtx.Core.CreateCategory(l.ctx, &core_client.CreateCategoryReq{
		Name: req.Name,
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
