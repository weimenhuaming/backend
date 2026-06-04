package agent

import (
	"context"
	"errors"
	"sync"

	"other-rpc/internal/agent/llm"
	"other-rpc/internal/agent/rag"
	"other-rpc/internal/agent/vector"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
)

// Agent 聚合 LLM、Embedding、检索器与 RAG 问答链。
type Agent struct {
	cfg      config.KnowledgeBaseConf
	brain    llms.Model
	embedder embeddings.Embedder

	// 检索器
	retriever vectorstores.Retriever
	qa        *rag.QA

	docCount   int
	chunkCount int
	mu         sync.RWMutex
}

// NewAgent 创建 Agent 运行时：仅加载已有检索器，不会在启动时重建。
func NewAgent(cfg config.KnowledgeBaseConf, em embeddings.Embedder, retriever vectorstores.Retriever) (*Agent, error) {
	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	a := &Agent{
		cfg:       cfg,
		brain:     chatModel,
		embedder:  em,
		retriever: retriever,
		qa:        rag.NewQA(chatModel, retriever),
	}
	return a, nil
}

// Build 从知识库目录构建向量索引并持久化
func (a *Agent) Build(ctx context.Context) (int, int, error) {
	cfg := a.cfg

	// 使用已有Embedder构建向量索引
	retriever, docCount, chunkCount, err := vector.BuildIndex(ctx, cfg, a.embedder)
	if err != nil {
		return 0, 0, err
	}

	a.docCount = docCount
	a.chunkCount = chunkCount
	a.retriever = retriever
	a.qa = rag.NewQA(a.brain, retriever)

	return docCount, chunkCount, nil
}

// Chat 基于已加载检索器的 RAG 问答。
func (a *Agent) Chat(ctx context.Context, question string) (string, error) {
	a.mu.RLock()
	qa := a.qa
	a.mu.RUnlock()
	if qa == nil {
		return "", errors.New("检索器未加载，请先构建知识库")
	}
	return qa.Ask(ctx, question)
}
