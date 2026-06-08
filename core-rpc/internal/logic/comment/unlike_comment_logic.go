package comment

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeCommentLogic {
	return &UnlikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnlikeCommentLogic) UnlikeComment(in *core.UnlikeCommentReq) (*core.UnlikeCommentResp, error) {
	if in.UserId == 0 || in.CommentId == 0 {
		return nil, errors.New("参数无效")
	}

	exists, err := l.svcCtx.CommentRepo.Exists(in.CommentId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("评论不存在")
	}

	likeCount, err := l.svcCtx.InteractionRepo.UnlikeComment(in.UserId, in.CommentId)
	if err != nil {
		return nil, err
	}
	return &core.UnlikeCommentResp{LikeCount: likeCount}, nil
}
