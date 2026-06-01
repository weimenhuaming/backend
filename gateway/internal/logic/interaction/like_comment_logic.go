package interaction

import (
	"context"

	core_client "core-rpc/core_client"
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

func (l *LikeCommentLogic) LikeComment(req *types.LikeCommentReq) (resp *types.LikeCommentResp, err error) {
	userID, ok, msg := likeUserFromCtx(l.ctx)
	if !ok {
		return &types.LikeCommentResp{Code: authFailCode(msg), Msg: msg}, nil
	}
	if req.CommentId == 0 {
		return &types.LikeCommentResp{Code: 400, Msg: "评论ID无效"}, nil
	}

	r, err := l.svcCtx.Core.LikeComment(l.ctx, &core_client.LikeCommentReq{
		CommentId: req.CommentId,
		UserId:    userID,
	})
	if err != nil {
		return &types.LikeCommentResp{Code: likeErrCode(err), Msg: err.Error()}, nil
	}

	return &types.LikeCommentResp{
		Code: 200,
		Msg:  "点赞成功",
		Data: types.LikeCommentData{LikeCount: r.GetLikeCount()},
	}, nil
}
