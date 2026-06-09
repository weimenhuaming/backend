// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"
	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CreateCommentData, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return nil, response.NewError(401, "用户未登录")
	}
	if req.ArticleId == 0 {
		return nil, response.NewError(400, "文章ID不存在")
	}
	if req.Content == "" {
		return nil, response.NewError(400, "评论内容不能为空")
	}

	r, err := l.svcCtx.Core.CreateComment(l.ctx, &core_client.CreateCommentReq{
		ArticleId: req.ArticleId,
		UserId:    userId,
		Content:   req.Content,
	})
	if err != nil {
		return nil, response.NewError(500, err.Error())
	}

	return &types.CreateCommentData{CommentId: r.GetCommentId()}, nil
}
