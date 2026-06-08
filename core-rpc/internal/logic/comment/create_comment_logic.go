package comment

import (
	"context"
	"errors"
	"strings"

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
	if in.UserId == 0 || in.ArticleId == 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(in.Content) == "" {
		return nil, errors.New("评论内容不能为空")
	}

	commentID, err := l.svcCtx.CommentRepo.CreateComment(in.UserId, in.ArticleId, strings.TrimSpace(in.Content))
	if err != nil {
		return nil, err
	}
	return &core.CreateCommentResp{CommentId: commentID}, nil
}
