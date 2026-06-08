package comment

import (
	"context"
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
	if err := l.svcCtx.CommentRepo.Delete(in.Id, in.UserId); err != nil {
		return nil, err
	}
	return &core.DeleteCommentResp{}, nil
}
