package logic

import (
	"context"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetCommentLikeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetCommentLikeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetCommentLikeStatusLogic {
	return &BatchGetCommentLikeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetCommentLikeStatusLogic) BatchGetCommentLikeStatus(in *core.BatchGetCommentLikeStatusReq) (*core.BatchGetCommentLikeStatusResp, error) {
	ids := uniquePositiveIDs(in.CommentIds)
	likedMap, err := l.svcCtx.InteractionRepo.BatchCommentLiked(l.svcCtx.Db, in.UserId, ids)
	if err != nil {
		return nil, err
	}

	items := make([]*core.CommentLikeStatusItem, 0, len(ids))
	for _, id := range ids {
		items = append(items, &core.CommentLikeStatusItem{
			CommentId: id,
			Liked:     likedMap[id],
		})
	}
	return &core.BatchGetCommentLikeStatusResp{Items: items}, nil
}

func uniquePositiveIDs(ids []uint64) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
