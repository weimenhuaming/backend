package agent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"other-rpc/internal/agent/embedding"
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
	retriever  vectorstores.Retriever
	qa         *rag.QA
	docCount   int
	chunkCount int
	mu         sync.RWMutex
}

// NewAgent 创建 Agent 运行时：仅加载已有检索器，不会在启动时重建。
func NewAgent(cfg config.KnowledgeBaseConf) (*Agent, error) {
	chatModel, err := llm.NewChatModel(cfg.LLM)
	if err != nil {
		return nil, err
	}

	embedder, err := embedding.NewEmbedder(cfg.Embedding)
	if err != nil {
		return nil, err
	}

	retriever, docCount, chunkCount, err := vector.Load(context.Background(), cfg, embedder)
	if err != nil {
		return nil, fmt.Errorf("加载检索器失败（请先调用 Build 构建知识库）: %w", err)
	}

	a := &Agent{
		cfg:      cfg,
		brain:    chatModel,
		embedder: embedder,
	}
	a.setRetriever(retriever, docCount, chunkCount)
	return a, nil
}

func (a *Agent) setRetriever(retriever vectorstores.Retriever, docCount, chunkCount int) {
	a.retriever = retriever
	a.docCount = docCount
	a.chunkCount = chunkCount
	a.qa = rag.NewQA(a.brain, retriever)
}

// Build 从知识库目录构建向量索引并持久化，同时热更新检索器。
func (a *Agent) Build(ctx context.Context, dataPath string) (docCount, chunkCount int, err error) {
	cfg := a.cfg
	if dataPath != "" {
		cfg.DataPath = dataPath
	}

	retriever, docCount, chunkCount, err := vector.BuildWithEmbedder(ctx, cfg, a.embedder)
	if err != nil {
		return 0, 0, err
	}

	a.mu.Lock()
	a.setRetriever(retriever, docCount, chunkCount)
	a.mu.Unlock()
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

// Stats 返回当前 Agent 的基础状态信息。
func (a *Agent) Stats() (docCount int, chunkCount int, topK int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	topK = a.cfg.TopK
	if topK <= 0 {
		topK = 4
	}
	if a.qa == nil {
		return 0, 0, topK
	}
	return a.docCount, a.chunkCount, topK
}
