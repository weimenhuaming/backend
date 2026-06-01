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

func (l *CreateReplyLogic) CreateReply(req *types.CreateReplyReq) (resp *types.CreateReplyResp, err error) {
	userId, ok := l.ctx.Value("X-user-Id").(uint64)
	if !ok || userId == 0 {
		return &types.CreateReplyResp{Code: 401, Msg: "用户未登录"}, nil
	}
	if req.RootId == 0 || req.ParentId == 0 {
		return &types.CreateReplyResp{Code: 400, Msg: "评论参数无效"}, nil
	}
	if req.Content == "" {
		return &types.CreateReplyResp{Code: 400, Msg: "回复内容不能为空"}, nil
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
		return &types.CreateReplyResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.CreateReplyResp{
		Code: 200,
		Msg:  "回复成功",
		Data: types.CreateReplyData{ReplyId: r.GetReplyId()},
	}, nil
}
