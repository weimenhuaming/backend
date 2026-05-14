package login

import (
	"context"
	"core-rpc/core"
	"gateway/internal/utils"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailLoginLogic) EmailLogin(req *types.LoginEmailReq) (resp *types.LoginResp, err error) {
	if !utils.IsValidEmail(req.Email) {
		return &types.LoginResp{
			Code: 400,
			Msg:  "邮箱格式不正确",
		}, nil
	}
	if req.Captcha == "" {
		return &types.LoginResp{
			Code: 400,
			Msg:  "验证码不能为空",
		}, nil
	}

	// 从缓存中获取验证码
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return &types.LoginResp{
			Code: 400,
			Msg:  "验证码不存在或者已过期",
		}, nil
	}
	if captcha != req.Captcha {
		return &types.LoginResp{
			Code: 400,
			Msg:  "验证码错误",
		}, nil
	}

	// 3.调用逻辑函数返回的是rpc中的返回值。
	RpcResp, err := l.svcCtx.Core.EmailLogin(l.ctx, &core.EmailLoginReq{
		Email: req.Email,
	})
	if err != nil {
		return &types.LoginResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	// 4.签发token
	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(RpcResp.Id, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return &types.LoginResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}
	refreshToken, err := jwt.GetRefreshToken(RpcResp.Id, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return &types.LoginResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	// 5.统一 API 响应
	return &types.LoginResp{
		Code: 200,
		Msg:  "登录成功",
		Data: types.LoginData{
			Id:           RpcResp.Id,
			Name:         RpcResp.Name,
			Phone:        RpcResp.Phone,
			Email:        RpcResp.Email,
			Avatar:       RpcResp.Avatar,
			Uuid:         RpcResp.Uuid,
			Role:         RpcResp.Role,
			Sex:          RpcResp.Sex,
			Age:          RpcResp.Age,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
