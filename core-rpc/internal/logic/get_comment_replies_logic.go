package logic

import (
	"context"
	"core-rpc/internal/utils"

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
	page := utils.NormalizePage(in.Page)
	size := utils.NormalizeSize(in.Size, 10)
	off, limit := utils.OffsetLimit(page, size)

	total, replies, err := l.svcCtx.CommentRepo.ListReplies(in.RootId, off, limit)
	if err != nil {
		return nil, err
	}

	// collect user ids
	ids := make([]uint64, 0, len(replies))
	seen := make(map[uint64]struct{})
	for _, r := range replies {
		if _, ok := seen[r.UserID]; !ok {
			seen[r.UserID] = struct{}{}
			ids = append(ids, r.UserID)
		}
	}
	userMap, err := l.svcCtx.UserRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	return &core.GetCommentRepliesResp{
		Replies: commentsToProtoList(replies, userMap, nil),
		Page:    int32(page),
		Size:    int32(size),
		Total:   int32(total),
	}, nil
}
