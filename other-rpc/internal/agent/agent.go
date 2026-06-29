package agent

import (
	"context"
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
	// 知识库配置
	cfg config.KnowledgeBaseConf

	// llm大脑
	brain llms.Model

	// 短期记忆
	memory *memory.RedisChatMessageHistory

	// 检索器
	retriever vectorstores.Retriever

	// RAG处理
	qa *rag.QA
}

// NewAgent 创建 Agent 运行时：加载 LLM 与已有检索器，不会在启动时重建索引。
func NewAgent(cfg config.KnowledgeBaseConf, retriever vectorstores.Retriever, cache *redis.Redis) (*Agent, error) {
	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	memoryInRedis := memory.NewRedisChatMessageHistory(cache, cfg.Memory.SessionTTL, cfg.Memory.WindowTurns)

	return &Agent{
		cfg:       cfg,
		brain:     chatModel,
		memory:    memoryInRedis,
		retriever: retriever,
		qa:        rag.NewQA(chatModel, retriever, memoryInRedis),
	}, nil
}

// SetRetriever 切换当前问答使用的检索器。(不同向量数据库)
func (a *Agent) SetRetriever(retriever vectorstores.Retriever) {
	a.retriever = retriever
	a.qa = rag.NewQA(a.brain, retriever, a.memory)
}

// Chat 只有 RAG，无 Memory,也就不需要sessionId
func (a *Agent) Chat(ctx context.Context, question string) (string, error) {
	return a.qa.Ask(ctx, question)
}

// ChatStream 流式 RAG 问答，检索完成后逐 token 回调 send。
func (a *Agent) ChatStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	return a.qa.AskStream(ctx, sessionID, question, send)
}
