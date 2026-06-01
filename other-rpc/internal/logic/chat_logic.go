package logic

import (
	"context"
	"errors"
	"strings"

	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatLogic) Chat(in *agent.ChatRequest) (*agent.ChatResponse, error) {
	if in == nil || strings.TrimSpace(in.Question) == "" {
		return nil, errors.New("问题不能为空")
	}

	answer, err := l.svcCtx.Agent.Chat(l.ctx, in.Question)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("知识库问答失败: %v", err)
		return nil, err
	}

	return &agent.ChatResponse{Answer: answer}, nil
}
