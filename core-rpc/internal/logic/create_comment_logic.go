package logic

import (
	"context"
	"errors"
	"fmt"

	"core-rpc/internal/model/comment"
	"core-rpc/internal/svc"
	"core-rpc/pb/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *core.CreateCommentReq) (*core.CreateCommentResp, error) {
	fmt.Println("here")
	if in.GetUserId() == 0 {
		return nil, errors.New("missing user id")
	}
	if in.GetArticleId() == 0 {
		return nil, errors.New("missing article id")
	}
	if in.GetContent() == "" {
		return nil, errors.New("content is required")
	}

	if _, err := l.svcCtx.ArticleModel.FindOneActive(l.ctx, in.GetArticleId()); err != nil {
		return nil, errors.New("article not found")
	}

	row := &comment.Comment{
		ArticleId: in.GetArticleId(),
		UserId:    in.GetUserId(),
		ParentId:  0,
		Content:   in.GetContent(),
	}
	result, err := l.svcCtx.CommentModel.Insert(l.ctx, row)
	if err != nil {
		return nil, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err = l.svcCtx.CommentModel.UpdateRootId(l.ctx, uint64(commentID)); err != nil {
		return nil, err
	}
	if err = l.svcCtx.ArticleModel.IncCommentCount(l.ctx, in.GetArticleId(), 1); err != nil {
		return nil, err
	}

	return &core.CreateCommentResp{CommentId: uint64(commentID)}, nil
}
