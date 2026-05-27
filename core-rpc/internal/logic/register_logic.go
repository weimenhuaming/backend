package logic

import (
	"context"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *core.RegisterReq) (*core.RegisterResp, error) {
	var count int64
	if err := l.svcCtx.Db.Model(&entity.User{}).Where("email = ?", in.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	u := &entity.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     "user",
		Sex:      "未知",
	}
	if err := l.svcCtx.Db.Create(u).Error; err != nil {
		return nil, err
	}
	return &core.RegisterResp{}, nil
}
