package user

import (
	"context"
	"strings"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/converter"
	"gateway/internal/utils/vaild"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(req *types.UpdateUserProfileReq) (resp *types.UserProfile, err error) {
	userId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return nil, response.ErrorUnauthorized("用户未登录")
	}

	rpcReq := &core_client.UpdateUserProfileReq{
		UserId: userId,
		Name:   strings.TrimSpace(req.Name),
		Phone:  strings.TrimSpace(req.Phone),
		Sex:    strings.TrimSpace(req.Sex),
		Age:    req.Age,
	}

	if avatar := strings.TrimSpace(req.Avatar); avatar != "" {
		if !vaild.IsAdmin(l.ctx) {
			return nil, response.ErrorForbidden("仅管理员可修改头像")
		}
		rpcReq.Avatar = avatar
	}

	r, err := l.svcCtx.Core.UpdateUserProfile(l.ctx, rpcReq)
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return converter.ToUserProfile(r.GetProfile()), nil
}
