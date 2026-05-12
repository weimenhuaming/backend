package login

import (
	"context"
	"gateway/internal/utils"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailLogic {
	return &SendEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendEmailLogic) SendEmail(req *types.EmailReq) (resp *types.EmailResp, err error) {
	// 1.拿到邮箱和生成验证码
	email := req.Email
	captcha, err := utils.GenerateCode()
	if err != nil {
		return nil, err
	}

	// 2.存入缓存
	if err = l.svcCtx.Cache.SetexCtx(
		l.ctx,
		email,
		captcha, // 值随意，存在即可
		60,
	); err != nil {
		return nil, err
	}

	if err = utils.SendEmailVerificationCode(email, captcha); err != nil {
		return nil, err
	}
	return &types.EmailResp{
		Success: true,
	}, nil
}
