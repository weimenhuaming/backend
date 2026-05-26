// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCommentsLogic {
	return &GetUserCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserCommentsLogic) GetUserComments(req *types.GetUserCommentsReq) (resp *types.GetUserCommentsResp, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return &types.GetUserCommentsResp{Code: 401, Msg: "用户未登录"}, nil
	}

	page := int32(req.Page)
	size := int32(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	r, err := l.svcCtx.Core.GetUserComments(l.ctx, &core_client.GetUserCommentsReq{
		UserId: userId,
		Page:   page,
		Size:   size,
	})
	if err != nil {
		return &types.GetUserCommentsResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.GetUserCommentsResp{
		Code: 200,
		Msg:  "ok",
		Data: types.GetUserCommentsData{
			Comments: toTypesCommentList(r.GetComments()),
			Total:    uint32(r.GetTotal()),
			Page:     uint32(r.GetPage()),
			PageSize: uint32(r.GetSize()),
		},
	}, nil
}
