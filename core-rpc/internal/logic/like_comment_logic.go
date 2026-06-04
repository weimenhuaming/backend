package logic

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeCommentLogic) LikeComment(in *core.LikeCommentReq) (*core.LikeCommentResp, error) {
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

	likeCount, err := l.svcCtx.InteractionRepo.LikeComment(in.UserId, in.CommentId)
	if err != nil {
		return nil, err
	}
	return &core.LikeCommentResp{LikeCount: likeCount}, nil
}
