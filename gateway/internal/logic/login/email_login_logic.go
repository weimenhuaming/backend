package login

import (
	"context"
	"core-rpc/core"
	"errors"
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
	// 1. 判断邮箱是否合理
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 2.从缓存中获取验证码
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return nil, errors.New("验证码不存在或者已过期")
	}
	if captcha != req.Captcha {
		return nil, errors.New("验证码错误")
	}

	// 3.调用逻辑函数返回的是rpc中的返回值。
	RpcResp, err := l.svcCtx.Core.EmailLogin(l.ctx, &core.EmailLoginReq{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	// 4.签发token
	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(RpcResp.Id, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.GetRefreshToken(RpcResp.Id, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return nil, err
	}

	// 5.统一 API 响应
	resp = &types.LoginResp{
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
	}
	return
}
