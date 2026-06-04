package logic

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(in *core.DeleteCategoryReq) (*core.DeleteCategoryResp, error) {
	rows, err := l.svcCtx.CateRepo.DeleteByID(in.Id)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("分类不存在")
	}
	return &core.DeleteCategoryResp{}, nil
}
