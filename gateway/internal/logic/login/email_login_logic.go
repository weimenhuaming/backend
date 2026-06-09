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

func (l *EmailLoginLogic) EmailLogin(req *types.LoginEmailReq) (resp *types.LoginData, err error) {
	if !utils.IsValidEmail(req.Email) {
		return nil, response.NewError(400, "邮箱格式不正确")
	}
	if req.Captcha == "" {
		return nil, response.NewError(400, "验证码不能为空")
	}

	// 从缓存中获取验证码
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return nil, response.NewError(400, "验证码不存在或者已过期")
	}
	if captcha != req.Captcha {
		return nil, response.NewError(400, "验证码错误")
	}

	// 3.调用逻辑函数返回的是rpc中的返回值。
	RpcResp, err := l.svcCtx.Core.EmailLogin(l.ctx, &core_client.EmailLoginReq{
		Email: req.Email,
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	// 4.签发token
	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(RpcResp.Id, RpcResp.Role, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}
	refreshToken, err := jwt.GetRefreshToken(RpcResp.Id, RpcResp.Role, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	// 5.统一 API 响应
	return &types.LoginData{
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
	}, nil
}
