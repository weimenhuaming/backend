// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentRepliesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentRepliesLogic {
	return &GetCommentRepliesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentRepliesLogic) GetCommentReplies(req *types.GetCommentRepliesReq) (resp *types.GetCommentRepliesData, err error) {
	if req.RootId == 0 {
		return nil, response.NewError(400, "根评论ID不存在")
	}

	page := int32(req.Page)
	size := int32(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	r, err := l.svcCtx.Core.GetCommentReplies(l.ctx, &core_client.GetCommentRepliesReq{
		RootId: req.RootId,
		Page:   page,
		Size:   size,
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	return &types.GetCommentRepliesData{
		Replies:  toTypesCommentList(r.GetReplies()),
		Total:    uint32(r.GetTotal()),
		Page:     uint32(r.GetPage()),
		PageSize: uint32(r.GetSize()),
	}, nil
}
