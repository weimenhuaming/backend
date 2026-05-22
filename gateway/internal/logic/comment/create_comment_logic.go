// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

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

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CreateCommentResp, err error) {
	// 从上下文中获取用户ID
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return &types.CreateCommentResp{
			Code: 400,
			Msg:  "用户ID不存在",
		}, nil
	}

	// 参数校验
	if req.ArticleId == 0 {
		return &types.CreateCommentResp{
			Code: 400,
			Msg:  "文章ID不存在",
		}, nil
	}
	if req.Content == "" {
		return &types.CreateCommentResp{
			Code: 400,
			Msg:  "评论内容不能为空",
		}, nil
	}

	return &types.CreateCommentResp{
		Code: 200,
		Msg:  "发表成功",
	}, nil
}
