// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package comment

import (
	"context"

	core_client "core-rpc/core_client"
	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	"gateway/internal/utils/vaild"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReplyLogic) CreateReply(req *types.CreateReplyReq) (resp *types.CreateReplyData, err error) {
	userId, ok := vaild.GetUserID(l.ctx)
	if !ok {
		return nil, response.ErrorUnauthorized("用户未登录")
	}
	if req.RootId == 0 || req.ParentId == 0 {
		return nil, response.ErrorBadRequest("评论参数无效")
	}
	if req.Content == "" {
		return nil, response.ErrorBadRequest("回复内容不能为空")
	}

	r, err := l.svcCtx.Core.CreateReply(l.ctx, &core_client.CreateReplyReq{
		RootId:      req.RootId,
		ParentId:    req.ParentId,
		UserId:      userId,
		ReplyToId:   req.ReplyToId,
		ReplyToName: req.ReplyToName,
		Content:     req.Content,
	})
	if err != nil {
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.CreateReplyData{ReplyId: r.GetReplyId()}, nil
}
