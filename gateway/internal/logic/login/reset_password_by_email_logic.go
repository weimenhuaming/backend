package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type Reset_password_by_emailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReset_password_by_emailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Reset_password_by_emailLogic {
	return &Reset_password_by_emailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Reset_password_by_emailLogic) Reset_password_by_email(req *types.ResetPasswordReq) (resp *types.ResetPasswordResp, err error) {
	if !utils.IsValidEmail(req.Email) {
		return &types.ResetPasswordResp{
			Code: 400,
			Msg:  "邮箱格式不正确",
		}, nil
	}
	if req.NewPassword == "" || req.Captcha == "" {
		return &types.ResetPasswordResp{
			Code: 400,
			Msg:  "密码或验证码不能为空",
		}, nil
	}
	if req.NewPassword != req.ConfirmPassword {
		return &types.ResetPasswordResp{
			Code: 400,
			Msg:  "两次输入的密码不一致",
		}, nil
	}

	key := "email:captcha:" + strings.ToLower(strings.TrimSpace(req.Email))
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, key)
	if err != nil {
		return &types.ResetPasswordResp{
			Code: 400,
			Msg:  "验证码不存在或已过期",
		}, nil
	}
	if captcha != req.Captcha {
		return &types.ResetPasswordResp{
			Code: 400,
			Msg:  "验证码错误",
		}, nil
	}

	_, err = l.svcCtx.Core.ResetPasswordByEmail(l.ctx, &core_client.ResetPasswordEmailReq{
		Email:    req.Email,
		Password: utils.Bcrypt(req.NewPassword),
	})
	if err != nil {
		return &types.ResetPasswordResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.ResetPasswordResp{
		Code: 200,
		Msg:  "密码重置成功",
	}, nil
}
