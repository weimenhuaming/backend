package login

import (
	"context"
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
	_, err := l.svcCtx.UserRepo.EmailRegister(in.Name, in.Email, in.Password)
	if err != nil {
		return nil, err
	}
	return &core.RegisterResp{}, nil
}
