package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCommentsLogic {
	return &GetUserCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserCommentsLogic) GetUserComments(in *core.GetUserCommentsReq) (*core.GetUserCommentsResp, error) {
	page, size := normalizePageSize(in.GetPage(), in.GetSize())

	rows, err := l.svcCtx.CommentModel.ListByUser(l.ctx, in.GetUserId(), page, size)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.CommentModel.CountByUser(l.ctx, in.GetUserId())
	if err != nil {
		return nil, err
	}

	return &core.GetUserCommentsResp{
		Comments: buildCommentInfoList(l.ctx, l.svcCtx, rows),
		Page:     page,
		Size:     size,
		Total:    int32(total),
	}, nil
}
