package logic

import (
	"context"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordByEmailLogic {
	return &ResetPasswordByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordByEmailLogic) ResetPasswordByEmail(in *core.ResetPasswordEmailReq) (*core.ResetPasswordEmailResp, error) {
	err := l.svcCtx.UserRepo.ResetPasswordByEmail(in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &core.ResetPasswordEmailResp{}, nil
}
