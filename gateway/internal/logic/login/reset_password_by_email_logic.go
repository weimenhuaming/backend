package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"

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

func (l *Reset_password_by_emailLogic) Reset_password_by_email(req *types.ResetPasswordReq) error {
	if !utils.IsValidEmail(req.Email) {
		return response.NewError(400, "邮箱格式不正确")
	}
	if req.NewPassword == "" || req.Captcha == "" {
		return response.NewError(400, "密码或验证码不能为空")
	}
	if req.NewPassword != req.ConfirmPassword {
		return response.NewError(400, "两次输入的密码不一致")
	}

	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return response.NewError(400, "验证码不存在或已过期")
	}
	if captcha != req.Captcha {
		return response.NewError(400, "验证码错误")
	}

	_, err = l.svcCtx.Core.ResetPasswordByEmail(l.ctx, &core_client.ResetPasswordEmailReq{
		Email:    req.Email,
		Password: utils.Bcrypt(req.NewPassword),
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
