package agent

import (
	"context"
	"gateway/internal/utils/vaild"
	"io"
	"strings"
	"time"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	agent_client "other-rpc/agent_client"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AgentChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAgentChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgentChatStreamLogic {
	return &AgentChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AgentChatStreamLogic) AgentChatStream(req *types.AgentChatReq, client chan<- *types.AgentChatStreamChunk) error {
	// 登入之后才可以提问。
	_, ok, msg := vaild.GetActionUserID(l.ctx)
	if !ok {
		return response.ErrorActionAuth(msg)
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		return response.ErrorBadRequest("问题不能为空")
	}

	// 表示客户端和服务端建立的会话id
	sessionId := strings.TrimSpace(req.SessionId)
	if sessionId == "" {
		sessionId = uuid.NewString()
	}
	messageId := uuid.NewString()

	// 设置超时，避免ai生成内容超时
	chatCtx, cancel := context.WithTimeout(context.WithoutCancel(l.ctx), 5*time.Minute)
	defer cancel()

	stream, err := l.svcCtx.Agent.ChatStream(chatCtx, &agent_client.ChatRequest{
		Question:  question,
		SessionId: sessionId,
	})
	if err != nil {
		l.Errorf("agent chat stream failed: %v", err)
		return response.ErrorInternalServer(err.Error())
	}

	// 处理stream
	for {
		chunk, err := stream.Recv()
		// EOF表示到文件末尾了（结尾标识符），就需要返回true告诉前端
		if err == io.EOF {
			client <- &types.AgentChatStreamChunk{
				Done:      true,
				SessionId: sessionId,
				MessageId: messageId,
			}
			return nil
		}
		if err != nil {
			l.Errorf("agent chat stream recv failed: %v", err)
			return response.ErrorInternalServer(err.Error())
		}

		// 检测到对话流结束时，封装一个“结束信号”并发送给客户端，然后退出处理流程
		if chunk.GetDone() {
			client <- &types.AgentChatStreamChunk{
				Done:      true,
				SessionId: sessionId,
				MessageId: messageId,
			}
			return nil
		}

		// 拿到对于的内容，没有则跳过
		content := chunk.GetContent()
		if content == "" {
			continue
		}

		// 写入这个chan中去，在handler中读取
		client <- &types.AgentChatStreamChunk{
			Content:   content,
			SessionId: sessionId,
			MessageId: messageId,
		}
	}
}
