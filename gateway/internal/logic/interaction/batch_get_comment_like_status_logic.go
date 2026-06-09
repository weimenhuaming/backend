package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetCommentLikeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchGetCommentLikeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetCommentLikeStatusLogic {
	return &BatchGetCommentLikeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchGetCommentLikeStatusLogic) BatchGetCommentLikeStatus(req *types.BatchGetCommentLikeStatusReq) (resp *types.BatchGetCommentLikeStatusData, err error) {
	if len(req.CommentIds) == 0 {
		return &types.BatchGetCommentLikeStatusData{Items: []types.CommentLikeStatusItem{}}, nil
	}

	r, err := l.svcCtx.Core.BatchGetCommentLikeStatus(l.ctx, &core_client.BatchGetCommentLikeStatusReq{
		CommentIds: req.CommentIds,
		UserId:     userIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	items := make([]types.CommentLikeStatusItem, 0, len(r.GetItems()))
	for _, item := range r.GetItems() {
		items = append(items, types.CommentLikeStatusItem{
			CommentId: item.GetCommentId(),
			Liked:     item.GetLiked(),
		})
	}

	return &types.BatchGetCommentLikeStatusData{Items: items}, nil
}
