package login

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/pb"

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

func (l *ResetPasswordByEmailLogic) ResetPasswordByEmail(in *pb.ResetPasswordEmailReq) (*pb.ResetPasswordEmailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ResetPasswordEmailResp{}, nil
}
