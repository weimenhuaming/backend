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

type ResetPasswordByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordByEmailLogic {
	return &ResetPasswordByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordByEmailLogic) ResetPasswordByEmail(in *core.ResetPasswordEmailReq) (*core.ResetPasswordEmailResp, error) {
	var u entity.User
	err := l.svcCtx.Db.Where("email = ?", in.Email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		logx.WithContext(l.ctx).Errorf("find user by email failed: %v", err)
		return nil, errors.New("查询用户失败")
	}

	if err := l.svcCtx.Db.Model(&u).Update("password", in.Password).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("update password failed, id=%d: %v", u.ID, err)
		return nil, errors.New("更新密码失败")
	}
	return &core.ResetPasswordEmailResp{}, nil
}
