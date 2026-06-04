package logic

import (
	"context"
	"errors"
	"strings"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCategoryLogic) CreateCategory(in *core.CreateCategoryReq) (*core.CreateCategoryResp, error) {
	// 除去两边的空格
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	// 这里会有个问题gorm就算失败也会导致id自增
	// 要么就是先查后增，或者就是允许空洞
	// check existing by name
	if _, err := l.svcCtx.CateRepo.FindByName(name); err == nil {
		return nil, errors.New("分类名称已存在")
	}
	// if not found, create
	if _, err := l.svcCtx.CateRepo.Create(name); err != nil {
		// try to detect duplicate error via string fallback
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("分类名称已存在")
		}
		return nil, err
	}
	return &core.CreateCategoryResp{}, nil
}
