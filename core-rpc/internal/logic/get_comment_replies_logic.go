package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentRepliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentRepliesLogic {
	return &GetCommentRepliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentRepliesLogic) GetCommentReplies(in *core.GetCommentRepliesReq) (*core.GetCommentRepliesResp, error) {
	page, size := normalizePageSize(in.GetPage(), in.GetSize())

	rows, err := l.svcCtx.CommentModel.ListRepliesByRoot(l.ctx, in.GetRootId(), page, size)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.CommentModel.CountRepliesByRoot(l.ctx, in.GetRootId())
	if err != nil {
		return nil, err
	}

	return &core.GetCommentRepliesResp{
		Replies: buildCommentInfoList(l.ctx, l.svcCtx, rows),
		Page:    page,
		Size:    size,
		Total:   int32(total),
	}, nil
}
