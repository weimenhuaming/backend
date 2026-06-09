package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"
	"gateway/internal/utils/vaild"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileLogic) GetUserProfile() (resp *types.UserProfile, err error) {
	userId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return nil, response.ErrorUnauthorized("用户未登录")
	}

	r, err := l.svcCtx.Core.GetUserProfile(l.ctx, &core_client.GetUserProfileReq{
		UserId: userId,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	profile := converter.ToUserProfile(r.GetProfile())
	return &profile, nil
}
