package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/utils"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) error {
	// 1. 参数与格式（BFF 层统一校验）
	if !utils.IsValidEmail(req.Email) {
		return response.NewError(400, "邮箱格式不正确")
	}
	if req.Name == "" || req.Password == "" || req.Captcha == "" {
		return response.NewError(400, "用户名、密码或验证码不能为空")
	}

	// 2. 判断验证码是否有效
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return response.NewError(400, "验证码不存在或已过期")
	}
	if captcha != req.Captcha {
		return response.NewError(400, "验证码错误")
	}

	// 3.获得rpc响应
	_, err = l.svcCtx.Core.Register(l.ctx, &core_client.RegisterReq{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.Bcrypt(req.Password),
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
