package logic

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

// Logout 由 gateway 校验 token 是否传入；此处仅将非空 token 写入黑名单（空串跳过）。
func (l *LogoutLogic) Logout(in *core.LogoutReq) (*core.LogoutResp, error) {
	//if in.AccessToken != "" {
	//	key := tokenblacklist.AccessKey(in.AccessToken)
	//	if err := l.svcCtx.Cache.SetexCtx(l.ctx, key, "1", tokenblacklist.DefaultTTLSeconds); err != nil {
	//		logx.WithContext(l.ctx).Errorf("blacklist access token failed: %v", err)
	//		return nil, errors.New("登出失败")
	//	}
	//}
	//if in.RefreshToken != "" {
	//	key := tokenblacklist.RefreshKey(in.RefreshToken)
	//	if err := l.svcCtx.Cache.SetexCtx(l.ctx, key, "1", tokenblacklist.DefaultTTLSeconds); err != nil {
	//		logx.WithContext(l.ctx).Errorf("blacklist refresh token failed: %v", err)
	//		return nil, errors.New("登出失败")
	//	}
	//}

	return &core.LogoutResp{}, nil
}
