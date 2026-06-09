package login

import (
	"context"
	"core-rpc/core_client"
	"gateway/internal/utils"
	"gateway/internal/utils/converter"
	"gateway/internal/utils/vaild"

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
	// 1. Verify if the email address complies with requirements
	if !vaild.IsValidEmail(req.Email) {
		return nil, response.ErrorBadRequest("邮箱格式不正确")
	}
	if req.Captcha == "" {
		return nil, response.ErrorBadRequest("验证码不能为空")
	}

	// 2.从缓存中获取验证码
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return nil, response.ErrorBadRequest("验证码不存在或者已过期")
	}
	if captcha != req.Captcha {
		return nil, response.ErrorBadRequest("验证码错误")
	}

	// 3.call rpc
	RpcResp, err := l.svcCtx.Core.EmailLogin(l.ctx, &core_client.EmailLoginReq{
		Email: req.Email,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	// 4.签发token
	jwt := utils.NewJWT(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.RefreshSecret)
	accessToken, err := jwt.GetAccessToken(RpcResp.Id, RpcResp.Role, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}
	refreshToken, err := jwt.GetRefreshToken(RpcResp.Id, RpcResp.Role, l.svcCtx.Config.RefreshExpire)
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return converter.ToLoginData(RpcResp, accessToken, refreshToken), nil
}
