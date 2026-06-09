package interaction

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeCommentLogic) LikeComment(req *types.LikeCommentReq) (resp *types.LikeCommentData, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return nil, authFailError(msg)
	}
	if req.CommentId == 0 {
		return nil, response.ErrorBadRequest("评论ID无效")
	}

	r, err := l.svcCtx.Core.LikeComment(l.ctx, &core_client.LikeCommentReq{
		CommentId: req.CommentId,
		UserId:    userID,
	})
	if err != nil {
		return nil, likeRPCError(err)
	}

	return &types.LikeCommentData{LikeCount: r.GetLikeCount()}, nil
}
