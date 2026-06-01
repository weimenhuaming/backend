package interaction

import (
	"context"

	core_client "core-rpc/core_client"
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

func (l *UnlikeCommentLogic) UnlikeComment(req *types.UnlikeCommentReq) (resp *types.UnlikeCommentResp, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return &types.UnlikeCommentResp{Code: authFailCode(msg), Msg: msg}, nil
	}
	if req.CommentId == 0 {
		return &types.UnlikeCommentResp{Code: 400, Msg: "评论ID无效"}, nil
	}

	r, err := l.svcCtx.Core.UnlikeComment(l.ctx, &core_client.UnlikeCommentReq{
		CommentId: req.CommentId,
		UserId:    userID,
	})
	if err != nil {
		return &types.UnlikeCommentResp{Code: likeErrCode(err), Msg: err.Error()}, nil
	}

	return &types.UnlikeCommentResp{
		Code: 200,
		Msg:  "已取消点赞",
		Data: types.LikeCommentData{LikeCount: r.GetLikeCount()},
	}, nil
}
