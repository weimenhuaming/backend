package logic

import (
	"context"

	"other-rpc/internal/svc"
	"other-rpc/pb/agent"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatStreamLogic) ChatStream(in *agent.ChatRequest, stream agent.Agent_ChatStreamServer) error {
	question := in.GetQuestion()

	err := l.svcCtx.Agent.ChatStream(l.ctx, question, func(chunk string) error {
		return stream.Send(&agent.ChatStreamChunk{Content: chunk})
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("知识库流式问答失败: %v", err)
		return status.Errorf(codes.Internal, "知识库流式问答失败: %v", err)
	}

	return stream.Send(&agent.ChatStreamChunk{Done: true})
}
