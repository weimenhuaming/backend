package logic

import (
	"context"
	"errors"

	"core-rpc/internal/model/comment"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateReplyLogic) CreateReply(in *core.CreateReplyReq) (*core.CreateReplyResp, error) {
	if in.GetUserId() == 0 {
		return nil, errors.New("missing user id")
	}
	if in.GetRootId() == 0 || in.GetParentId() == 0 {
		return nil, errors.New("invalid root or parent id")
	}
	if in.GetContent() == "" {
		return nil, errors.New("content is required")
	}

	root, err := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetRootId())
	if err != nil {
		return nil, errors.New("root comment not found")
	}
	if root.ParentId != 0 {
		return nil, errors.New("root id must be a top-level comment")
	}

	parent, err := l.svcCtx.CommentModel.FindOneActive(l.ctx, in.GetParentId())
	if err != nil {
		return nil, errors.New("parent comment not found")
	}
	if parent.RootId != in.GetRootId() && parent.Id != in.GetRootId() {
		return nil, errors.New("parent comment does not belong to this thread")
	}

	row := &comment.Comment{
		ArticleId:   root.ArticleId,
		UserId:      in.GetUserId(),
		ParentId:    in.GetParentId(),
		RootId:      in.GetRootId(),
		ReplyToId:   in.GetReplyToId(),
		ReplyToName: in.GetReplyToName(),
		Content:     in.GetContent(),
	}
	result, err := l.svcCtx.CommentModel.Insert(l.ctx, row)
	if err != nil {
		return nil, err
	}
	replyID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err = l.svcCtx.CommentModel.IncChildCount(l.ctx, in.GetRootId(), 1); err != nil {
		return nil, err
	}
	if err = l.svcCtx.ArticleModel.IncCommentCount(l.ctx, root.ArticleId, 1); err != nil {
		return nil, err
	}

	return &core.CreateReplyResp{ReplyId: uint64(replyID)}, nil
}
