package logic

import (
	"context"
	"errors"

	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(in *core.DeleteCommentReq) (*core.DeleteCommentResp, error) {
	uid := in.GetUserId()
	if uid == 0 {
		return nil, errors.New("missing user id")
	}
	if in.GetId() == 0 {
		return nil, errors.New("missing comment id")
	}

	c, err := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetId())
	if err != nil {
		return nil, errors.New("comment not found")
	}
	if c.UserId != uid {
		return nil, errors.New("not comment owner")
	}

	if err = l.svcCtx.CommentModel.SoftDelete(l.ctx, in.GetId(), uid); err != nil {
		return nil, err
	}

	if c.ParentId > 0 && c.RootId > 0 {
		_ = l.svcCtx.CommentModel.IncChildCount(l.ctx, c.RootId, -1)
	}
	_ = l.svcCtx.ArticleModel.IncCommentCount(l.ctx, c.ArticleId, -1)

	return &core.DeleteCommentResp{}, nil
}
