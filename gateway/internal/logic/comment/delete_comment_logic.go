// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	core_client "core-rpc/core_client"
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

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentResp, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return &types.DeleteCommentResp{Code: 401, Msg: "用户未登录"}, nil
	}
	if req.Id == 0 {
		return &types.DeleteCommentResp{Code: 400, Msg: "评论ID不存在"}, nil
	}

	_, err = l.svcCtx.Core.DeleteComment(l.ctx, &core_client.DeleteCommentReq{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return &types.DeleteCommentResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.DeleteCommentResp{Code: 200, Msg: "删除成功"}, nil
}
