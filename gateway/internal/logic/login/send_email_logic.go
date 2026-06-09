package login

import (
	"context"
	"gateway/internal/utils"

	"gateway/internal/response"
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

func (l *SendEmailLogic) SendEmail(req *types.EmailReq) error {
	email := req.Email
	if !utils.IsValidEmail(email) {
		return response.NewError(400, "邮箱格式不正确")
	}

	captcha, err := utils.GenerateCode()
	if err != nil {
		return response.NewError(500, err.Error())
	}

	if err = l.svcCtx.Cache.SetexCtx(l.ctx, email, captcha, 60); err != nil {
		return response.NewError(500, err.Error())
	}

	if err = utils.SendEmailVerificationCode(email, captcha); err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
