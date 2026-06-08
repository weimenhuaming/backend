package comment

import (
	"context"
	"errors"
	"strings"

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
	if in.UserId == 0 || in.RootId == 0 || in.ParentId == 0 {
		return nil, errors.New("参数无效")
	}
	if strings.TrimSpace(in.Content) == "" {
		return nil, errors.New("回复内容不能为空")
	}

	replyID, err := l.svcCtx.CommentRepo.CreateReply(in.UserId, in.RootId, in.ParentId, in.ReplyToId, in.ReplyToName, strings.TrimSpace(in.Content))
	if err != nil {
		return nil, err
	}
	return &core.CreateReplyResp{ReplyId: replyID}, nil
}
