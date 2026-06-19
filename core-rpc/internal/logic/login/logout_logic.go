package login

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *core.LogoutReq) (*core.LogoutResp, error) {
	LogoutKey := "blacklist:" + in.RefreshToken
	err := l.svcCtx.Cache.Setex(LogoutKey, "1", 604800)
	if err != nil {
		return nil, err
	}

	return &core.LogoutResp{}, nil
}
