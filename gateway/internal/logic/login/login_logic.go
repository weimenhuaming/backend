package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/response"
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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginData, err error) {
	name := req.Name
	password := req.Password
	if name == "" || password == "" {
		return nil, response.ErrorBadRequest("用户名和密码不能为空")
	}
	if req.CaptchaId == "" || req.Code == "" {
		return nil, response.ErrorBadRequest("请填写验证码")
	}

	if !base64Captcha.DefaultMemStore.Verify(req.CaptchaId, req.Code, true) {
		return nil, response.ErrorBadRequest("验证码错误或已过期")
	}

	rpcResp, err := l.svcCtx.Core.NameLogin(l.ctx, &core_client.NameLoginReq{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return nil, response.ErrorBadRequest(err.Error())
	}

	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(rpcResp.Id, rpcResp.Role, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}
	refreshToken, err := jwt.GetRefreshToken(rpcResp.Id, rpcResp.Role, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.LoginData{
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
	}, nil
}
