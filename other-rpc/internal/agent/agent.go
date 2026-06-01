package agent

import (
	"context"

	"other-rpc/internal/agent/embedding"
	"other-rpc/internal/agent/kb"
	"other-rpc/internal/agent/llm"
	"other-rpc/internal/agent/rag"
	"other-rpc/internal/config"
)

// Agent 是最终对外暴露的运行时模式，聚合 LLM / KB / RAG。
type Agent struct {
	index *kb.Index
	qa    *rag.QA
}

// New 创建完整的 Agent 运行时。
func New(cfg config.KnowledgeBaseConf) (*Agent, error) {
	ctx := context.Background()

	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	embedder, err := embedding.NewEmbedder(cfg.Embedding)
	if err != nil {
		return nil, err
	}

	index, err := kb.Build(ctx, cfg, embedder)
	if err != nil {
		return nil, err
	}

	return &Agent{
		index: index,
		qa:    rag.NewQA(chatModel, index.Retriever),
	}, nil
}

// Chat 对话接口（当前实现为 RAG 问答）。
func (a *Agent) Chat(ctx context.Context, question string) (string, error) {
	return a.qa.Ask(ctx, question)
}

// Stats 返回当前 Agent 的基础状态信息。
func (a *Agent) Stats() (docCount int, chunkCount int, topK int) {
	return a.index.DocCount, a.index.ChunkCount, a.index.TopK
}
