package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	name := req.Name
	password := req.Password
	if name == "" || password == "" {
		return &types.LoginResp{Code: 400, Msg: "用户名和密码不能为空"}, nil
	}
	if req.CaptchaId == "" || req.Code == "" {
		return &types.LoginResp{Code: 400, Msg: "请填写验证码"}, nil
	}

	if !base64Captcha.DefaultMemStore.Verify(req.CaptchaId, req.Code, true) {
		return &types.LoginResp{Code: 400, Msg: "验证码错误或已过期"}, nil
	}

	rpcResp, err := l.svcCtx.Core.NameLogin(l.ctx, &core_client.NameLoginReq{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return &types.LoginResp{Code: 400, Msg: err.Error()}, nil
	}

	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(rpcResp.Id, rpcResp.Role, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return &types.LoginResp{Code: 500, Msg: err.Error()}, nil
	}
	refreshToken, err := jwt.GetRefreshToken(rpcResp.Id, rpcResp.Role, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return &types.LoginResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.LoginResp{
		Code: 200,
		Msg:  "登录成功",
		Data: types.LoginData{
			Id:           rpcResp.Id,
			Name:         rpcResp.Name,
			Phone:        rpcResp.Phone,
			Email:        rpcResp.Email,
			Avatar:       rpcResp.Avatar,
			Uuid:         rpcResp.Uuid,
			Role:         rpcResp.Role,
			Sex:          rpcResp.Sex,
			Age:          rpcResp.Age,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
