package agent

import (
	"context"
	"strings"
	"time"

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

func (l *AgentChatLogic) AgentChat(req *types.AgentChatReq) (resp *types.AgentChatResp, err error) {
	question := strings.TrimSpace(req.Question)
	if question == "" {
		return &types.AgentChatResp{Code: 400, Msg: "问题不能为空"}, nil
	}

	sessionId := strings.TrimSpace(req.SessionId)
	if sessionId == "" {
		sessionId = uuid.NewString()
	}

	r, err := l.svcCtx.Agent.Chat(l.ctx, &agent_client.ChatRequest{
		Question: question,
	})
	if err != nil {
		l.Errorf("agent chat failed: %v", err)
		return &types.AgentChatResp{Code: 500, Msg: err.Error()}, nil
	}

	return &types.AgentChatResp{
		Code: 200,
		Msg:  "ok",
		Data: types.AgentChatData{
			SessionId: sessionId,
			Answer:    r.GetAnswer(),
			MessageId: uuid.NewString(),
			Timestamp: time.Now().Unix(),
		},
	}, nil
}
