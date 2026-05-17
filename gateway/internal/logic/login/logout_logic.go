package login

import (
	"context"
	"core-rpc/core_client"

	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	_, err = l.svcCtx.Core.Logout(l.ctx, &core_client.LogoutReq{})
	if err != nil {
		return &types.LogoutResp{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}
	return &types.LogoutResp{
		Code: 200,
		Msg:  "success",
	}, nil
}
