package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/utils"
	"strings"

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

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 参数与格式（BFF 层统一校验）
	if !utils.IsValidEmail(req.Email) {
		return &types.RegisterResp{
			Code: 400,
			Msg:  "邮箱格式不正确",
		}, nil
	}
	if req.Name == "" || req.Password == "" || req.Captcha == "" {
		return &types.RegisterResp{
			Code: 400,
			Msg:  "用户名、密码或验证码不能为空",
		}, nil
	}

	// 2. 判断验证码是否有效（使用与发送时相同的 key 前缀与规范化）
	key := "email:captcha:" + strings.ToLower(strings.TrimSpace(req.Email))
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, key)
	if err != nil {
		return &types.RegisterResp{
			Code: 400,
			Msg:  "验证码不存在或已过期",
		}, nil
	}
	if captcha != req.Captcha {
		return &types.RegisterResp{
			Code: 400,
			Msg:  "验证码错误",
		}, nil
	}

	// 3.获得rpc响应
	_, err = l.svcCtx.Core.Register(l.ctx, &core_client.RegisterReq{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.Bcrypt(req.Password),
	})
	if err != nil {
		return &types.RegisterResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &types.RegisterResp{
		Code: 200,
		Msg:  "注册成功",
	}, nil
}
