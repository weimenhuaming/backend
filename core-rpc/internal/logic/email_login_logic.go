package logic

import (
	"context"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmailLoginLogic) EmailLogin(in *core.EmailLoginReq) (*core.LoginResp, error) {
	u, err := l.svcCtx.UserRepo.EmailLogin(in.Email)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("email login failed: %v", err)
		return nil, err
	}
	if u == nil {
		logx.WithContext(l.ctx).Errorf("email login failed: email not exist")
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
