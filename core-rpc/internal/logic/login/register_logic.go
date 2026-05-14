package login

import (
	"context"
	"core-rpc/internal/model/user"
	"core-rpc/internal/svc"
	"core-rpc/pb/pb"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 检查邮箱是否已存在
	_, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil {
		return nil, errors.New("邮箱已存在")
	}
	if !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}

	// 创建用户
	newUser := &user.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Role:     "user",
		Sex:      "unknown",
		Age:      0,
		Phone:    "",
		Avatar:   "",
	}

	_, err = l.svcCtx.UserModel.Insert(l.ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResp{}, nil
}
