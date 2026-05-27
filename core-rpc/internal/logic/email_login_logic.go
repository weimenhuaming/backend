package logic

import (
	"context"
	"errors"

	"core-rpc/internal/model/entity"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	var u entity.User
	err := l.svcCtx.Db.Where("email = ?", in.Email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		logx.WithContext(l.ctx).Errorf("find user by email failed: %v", err)
		return nil, errors.New("查询用户失败")
	}

	return &core.LoginResp{
		Id:     u.ID,
		Name:   u.Name,
		Phone:  u.Phone,
		Email:  u.Email,
		Avatar: u.Avatar,
		Role:   u.Role,
		Sex:    u.Sex,
		Age:    u.Age,
	}, nil
}
