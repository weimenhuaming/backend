package agent

import (
	"context"
	"errors"
	"other-rpc/internal/agent/memory"
	"sync"

	"other-rpc/internal/agent/llm"
	"other-rpc/internal/agent/rag"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
)

// Agent 聚合 LLM、检索器与 RAG 问答链，仅负责对话。
type Agent struct {
	cfg      config.KnowledgeBaseConf
	brain    llms.Model
	sessions *memory.Store

	mu        sync.RWMutex
	retriever vectorstores.Retriever
	qa        *rag.QA
}

// NewAgent 创建 Agent 运行时：加载 LLM 与已有检索器，不会在启动时重建索引。
func NewAgent(cfg config.KnowledgeBaseConf, retriever vectorstores.Retriever) (*Agent, error) {
	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	sessions := memory.NewStore(cfg.Memory)

	a := &Agent{
		cfg:       cfg,
		brain:     chatModel,
		sessions:  sessions,
		retriever: retriever,
	}
	a.qa = rag.NewQA(chatModel, retriever, sessions)
	return a, nil
}

// SetRetriever 切换当前问答使用的检索器。
func (a *Agent) SetRetriever(retriever vectorstores.Retriever) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.retriever = retriever
	a.qa = rag.NewQA(a.brain, retriever, a.sessions)
}

// Chat 基于已加载检索器的 RAG 问答，sessionID 用于短期多轮记忆。
func (a *Agent) Chat(ctx context.Context, sessionID, question string) (string, error) {
	a.mu.RLock()
	qa := a.qa
	a.mu.RUnlock()
	if qa == nil {
		return "", errors.New("检索器未加载，请先构建知识库")
	}
	return qa.Ask(ctx, sessionID, question)
}

// ChatStream 流式 RAG 问答，检索完成后逐 token 回调 send。
func (a *Agent) ChatStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	a.mu.RLock()
	qa := a.qa
	a.mu.RUnlock()
	if qa == nil {
		return errors.New("检索器未加载，请先构建知识库")
	}
	return qa.AskStream(ctx, sessionID, question, send)
}
