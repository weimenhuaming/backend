package user

import (
	"context"
	"core-rpc/internal/model/converter"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserProfileLogic) GetUserProfile(in *core.GetUserProfileReq) (*core.GetUserProfileResp, error) {
	if in.GetUserId() == 0 {
		return nil, errors.New("用户 ID 无效")
	}

	user, err := l.svcCtx.UserRepo.FindProfileByID(in.GetUserId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &core.GetUserProfileResp{
		Profile: converter.UserToProfile(user),
	}, nil
}
