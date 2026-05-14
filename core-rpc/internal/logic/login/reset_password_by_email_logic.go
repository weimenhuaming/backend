package login

import (
	"context"
	"core-rpc/internal/model/user"
	"core-rpc/internal/svc"
	"core-rpc/pb/pb"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
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

// ResetPasswordByEmail 由 gateway 完成参数判空、验证码与 Redis 校验；此处仅做领域内的用户查找与密码更新。
func (l *ResetPasswordByEmailLogic) ResetPasswordByEmail(in *pb.ResetPasswordEmailReq) (*pb.ResetPasswordEmailResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, errors.New("该邮箱尚未注册")
		}
		logx.WithContext(l.ctx).Errorf("FindOneByEmail failed: %v", err)
		return nil, errors.New("查询用户失败")
	}

	u.Password = in.Password
	if err := l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		logx.WithContext(l.ctx).Errorf("update password failed, id=%d: %v", u.Id, err)
		return nil, errors.New("更新密码失败")
	}

	return &pb.ResetPasswordEmailResp{}, nil
}
