package agent

import (
	"context"
	"errors"
	"sync"

	"other-rpc/internal/agent/llm"
	"other-rpc/internal/agent/memory"
	"other-rpc/internal/agent/rag"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Agent 聚合 LLM、检索器与 RAG 问答链，仅负责对话。
type Agent struct {
	cfg   config.KnowledgeBaseConf
	brain llms.Model
	cache *redis.Redis

	mu        sync.RWMutex
	retriever vectorstores.Retriever
	qa        *rag.QA
}

// NewAgent 创建 Agent 运行时：加载 LLM 与已有检索器，不会在启动时重建索引。
func NewAgent(cfg config.KnowledgeBaseConf, retriever vectorstores.Retriever, cache *redis.Redis) (*Agent, error) {
	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	a := &Agent{
		cfg:       cfg,
		brain:     chatModel,
		cache:     cache,
		retriever: retriever,
	}
	a.qa = rag.NewQA(chatModel, retriever, nil, cfg.Memory.WindowTurns)
	return a, nil
}

// SetRetriever 切换当前问答使用的检索器。
func (a *Agent) SetRetriever(retriever vectorstores.Retriever) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.retriever = retriever
	a.qa = rag.NewQA(a.brain, retriever, nil, a.cfg.Memory.WindowTurns)
}

// Chat 只有 RAG，无 Memory。
func (a *Agent) Chat(ctx context.Context, sessionID, question string) (string, error) {
	a.mu.RLock()
	qa := a.qa
	a.mu.RUnlock()
	if qa == nil {
		return "", errors.New("检索器未加载，请先构建知识库")
	}
	return qa.Ask(ctx, question)
}

// ChatStream 流式 RAG 问答，检索完成后逐 token 回调 send。
func (a *Agent) ChatStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	a.mu.RLock()
	brain := a.brain
	retriever := a.retriever
	memoryCfg := a.cfg.Memory
	cache := a.cache
	a.mu.RUnlock()

	if retriever == nil {
		return errors.New("检索器未加载，请先构建知识库")
	}

	chatHistory := memory.NewRedisChatMessageHistory(cache, sessionID, memoryCfg.SessionTTL)
	qa := rag.NewQA(brain, retriever, chatHistory, memoryCfg.WindowTurns)
	return qa.AskStream(ctx, question, send)
}
