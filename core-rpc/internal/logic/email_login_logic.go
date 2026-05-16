package logic

import (
	"context"
	"core-rpc/internal/model/user"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EmailLoginLogic) EmailLogin(in *core.EmailLoginReq) (*core.LoginResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		logx.WithContext(l.ctx).Errorf("FindOneByEmail failed, email=%s, err=%v", in.Email, err)
		return nil, errors.New("查询用户失败")
	}

	// token 由 gateway 层统一签发，这里只负责返回用户信息
	return &core.LoginResp{
		Id:     u.Id,
		Name:   u.Name,
		Phone:  u.Phone,
		Email:  u.Email,
		Avatar: u.Avatar,
		Role:   u.Role,
		Sex:    u.Sex,
		Age:    u.Age,
	}, nil
}
