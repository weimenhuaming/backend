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
	// 后续如果用户量大，存数据库的黑名单表就行，目前先存缓存
	LogoutKey := "blacklist:" + in.RefreshToken
	err := l.svcCtx.Cache.Setex(LogoutKey, "1", 604800)
	if err != nil {
		return nil, err
	}

	// 后续存入数据库
	return &core.LogoutResp{}, nil
}
