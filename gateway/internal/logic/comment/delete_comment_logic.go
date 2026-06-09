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

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) error {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return response.NewError(401, "用户未登录")
	}
	if req.Id == 0 {
		return response.NewError(400, "评论ID不存在")
	}

	_, err := l.svcCtx.Core.DeleteComment(l.ctx, &core_client.DeleteCommentReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return response.NewError(500, err.Error())
	}

	return nil
}
