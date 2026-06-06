package login

import (
	"context"
	"gateway/internal/utils"
	"strings"

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
	email := req.Email
	if !utils.IsValidEmail(email) {
		return &types.EmailResp{
			Code: 400,
			Msg:  "邮箱格式不正确",
		}, nil
	}

	// 拿到邮箱和生成验证码
	captcha, err := utils.GenerateCode()
	if err != nil {
		return &types.EmailResp{
			Code: 500,
			Msg:  err.Error(),
		}, err
	}

	// 2. 存入缓存（使用统一前缀并规范化邮箱，避免大小写/空格导致匹配失败）
	key := "email:captcha:" + strings.ToLower(strings.TrimSpace(email))
	if err = l.svcCtx.Cache.SetexCtx(
		l.ctx,
		key,
		captcha,
		300,
	); err != nil {
		return &types.EmailResp{
			Code: 500,
			Msg:  err.Error(),
		}, err
	}

	if err = utils.SendEmailVerificationCode(email, captcha); err != nil {
		return &types.EmailResp{
			Code: 500,
			Msg:  err.Error(),
		}, err
	}
	return &types.EmailResp{
		Code: 200,
		Msg:  "发送成功",
	}, nil
}
