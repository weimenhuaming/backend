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
	// 1. 判断邮箱是否合理
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 2. 判断验证码是否有效
	captcha, err := l.svcCtx.Cache.GetCtx(l.ctx, req.Email)
	if err != nil {
		return &types.RegisterResp{
			Code: 500,
			Msg:  errors.New("验证码失效").Error(),
		}, err
	}
	if captcha != req.Captcha {
		return &types.RegisterResp{
			Code: 500,
			Msg:  errors.New("验证码错误").Error(),
		}, err
	}

	// 2.获得rpc响应
	_, err = l.svcCtx.Core.Register(l.ctx, &core.RegisterReq{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.Bcrypt(req.Password),
	})
	if err != nil {
		return &types.RegisterResp{
			Code: 500,
			Msg:  err.Error(),
		}, err
	}

	//3.处理响应
	return &types.RegisterResp{
		Code: 200,
		Msg:  "注册成功",
	}, nil
}
