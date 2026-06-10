package agent

import (
	"context"
	"gateway/internal/utils/vaild"
	"strings"
	"time"

	"gateway/internal/response"
	"gateway/internal/svc"
	"gateway/internal/types"
	agent_client "other-rpc/agent_client"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AgentChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAgentChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgentChatLogic {
	return &AgentChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AgentChatLogic) AgentChat(req *types.AgentChatReq) (resp *types.AgentChatData, err error) {
	_, ok, msg := vaild.GetActionUserID(l.ctx)
	if !ok {
		return nil, response.ErrorActionAuth(msg)
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		return nil, response.ErrorBadRequest("问题不能为空")
	}

	sessionId := strings.TrimSpace(req.SessionId)
	if sessionId == "" {
		sessionId = uuid.NewString()
	}

	// Agent 问答包含向量检索 + LLM 生成，使用独立超时避免继承过短的 HTTP 上下文。
	chatCtx, cancel := context.WithTimeout(context.WithoutCancel(l.ctx), 2*time.Minute)
	defer cancel()

	r, err := l.svcCtx.Agent.Chat(chatCtx, &agent_client.ChatRequest{
		Question:  question,
		SessionId: sessionId,
	})
	if err != nil {
		l.Errorf("agent chat failed: %v", err)
		return nil, response.ErrorInternalServer(err.Error())
	}

	return &types.AgentChatData{
		SessionId: sessionId,
		Answer:    r.GetAnswer(),
		MessageId: uuid.NewString(),
		Timestamp: time.Now().Unix(),
	}, nil
}
