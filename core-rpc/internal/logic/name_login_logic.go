package logic

import (
	"context"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type NameLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNameLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NameLoginLogic {
	return &NameLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NameLoginLogic) NameLogin(in *core.NameLoginReq) (*core.LoginResp, error) {
	u, err := l.svcCtx.UserRepo.NameLogin(in.Name, in.Password)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("name login failed: %v", err)
		return nil, err
	}

	return &core.LoginResp{
		Id:     u.ID,
		Name:   u.Name,
		Phone:  u.Phone,
		Email:  u.Email,
		Avatar: u.Avatar,
		Role:   u.Role,
		Sex:    u.Sex,
		Age:    u.Age,
	}, nil
}
