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
	if in.GetUserId() == 0 || in.GetCommentId() == 0 {
		return nil, errors.New("invalid request")
	}

	if _, err := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetCommentId()); err != nil {
		return nil, errors.New("comment not found")
	}

	liked, err := l.svcCtx.CommentModel.FindCommentLikeActive(l.ctx, in.GetUserId(), in.GetCommentId())
	if err != nil {
		return nil, err
	}

	var delta int64
	if in.GetIsLike() {
		if liked {
			c, _ := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetCommentId())
			if c != nil {
				return &core.LikeCommentResp{LikeCount: uint32(c.LikeCount)}, nil
			}
			return &core.LikeCommentResp{}, nil
		}
		if err = l.svcCtx.CommentModel.InsertCommentLike(l.ctx, in.GetUserId(), in.GetCommentId()); err != nil {
			if restoreErr := l.svcCtx.CommentModel.RestoreCommentLike(l.ctx, in.GetUserId(), in.GetCommentId()); restoreErr != nil {
				return nil, err
			}
		}
		delta = 1
	} else {
		if !liked {
			c, _ := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetCommentId())
			if c != nil {
				return &core.LikeCommentResp{LikeCount: uint32(c.LikeCount)}, nil
			}
			return &core.LikeCommentResp{}, nil
		}
		if err = l.svcCtx.CommentModel.SoftDeleteCommentLike(l.ctx, in.GetUserId(), in.GetCommentId()); err != nil {
			return nil, err
		}
		delta = -1
	}

	likeCount, err := l.svcCtx.CommentModel.IncLikeCount(l.ctx, in.GetCommentId(), delta)
	if err != nil {
		return nil, err
	}

	return &core.LikeCommentResp{LikeCount: uint32(likeCount)}, nil
}
