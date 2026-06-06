package user

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"

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

func (l *GetUserProfileLogic) GetUserProfile() (resp *types.GetUserProfileResp, err error) {
	userId, ok := currentUserID(l.ctx)
	if !ok {
		return &types.GetUserProfileResp{Code: 401, Msg: "用户未登录"}, nil
	}

	r, err := l.svcCtx.Core.GetUserProfile(l.ctx, &core_client.GetUserProfileReq{
		UserId: userId,
	})
	if err != nil {
		return &types.GetUserProfileResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.GetUserProfileResp{
		Code: 200,
		Msg:  "ok",
		Data: toTypesUserProfile(r.GetProfile()),
	}, nil
}
