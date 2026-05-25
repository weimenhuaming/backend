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
	if in == nil || in.Id == 0 {
		return nil, errors.New("invalid category id")
	}

	// 检查该分类下是否存在未删除的文章
	if l.svcCtx.ArticleModel != nil {
		cnt, err := l.svcCtx.ArticleModel.CountByCategory(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("cannot delete category: articles exist in this category")
		}
	}

	if err := l.svcCtx.CategoryModel.Delete(l.ctx, in.Id); err != nil {
		return nil, err
	}

	return &core.DeleteCategoryResp{}, nil
}
