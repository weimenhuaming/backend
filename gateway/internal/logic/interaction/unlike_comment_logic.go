package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeCommentLogic {
	return &UnlikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeCommentLogic) UnlikeComment(req *types.UnlikeCommentReq) (resp *types.LikeCommentData, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return nil, authFailError(msg)
	}
	if req.CommentId == 0 {
		return nil, response.ErrorBadRequest("评论ID无效")
	}

	r, err := l.svcCtx.Core.UnlikeComment(l.ctx, &core_client.UnlikeCommentReq{
		CommentId: req.CommentId,
		UserId:    userID,
	})
	if err != nil {
		return nil, likeRPCError(err)
	}

	return &types.LikeCommentData{LikeCount: r.GetLikeCount()}, nil
}
