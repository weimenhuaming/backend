package login

import (
	"context"
	"core-rpc/core"
	"errors"
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
	// 1.首先判断验证码是不是错的
	if ok := base64Captcha.DefaultMemStore.Verify(req.CaptchaId, req.Code, true); !ok {
		return nil, errors.New("验证码错误")
	}

	// 2.调用逻辑函数返回的是rpc中的返回值。
	RpcResp, err := l.svcCtx.Core.Login(l.ctx, &core.LoginReq{
		Name:     req.Name,
		Password: utils.Bcrypt(req.Password),
	})
	if err != nil {
		return nil, err
	}

	// 3.处理完之后返回即可，把rpc的Resp给api。
	resp = &types.LoginResp{
		Id:     RpcResp.Id,
		Name:   RpcResp.Name,
		Phone:  RpcResp.Phone,
		Email:  RpcResp.Email,
		Avatar: RpcResp.Avatar,
		Uuid:   RpcResp.Uuid,
		Role:   RpcResp.Role,
		Age:    RpcResp.Age,
		Sex:    RpcResp.Sex,
	}
	// 4. 鉴权别忘记了
	// todo

	return
}
