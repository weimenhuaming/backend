package login

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCaptchaLogic {
	return &EmailCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmailCaptchaLogic) EmailCaptcha(in *pb.EmailCaptchaReq) (*pb.EmailCaptchaResp, error) {
	// todo: add your logic here and delete this line

	return &pb.EmailCaptchaResp{}, nil
}
