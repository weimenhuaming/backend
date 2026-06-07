package svc

import (
	"context"
	"fmt"
	"other-rpc/internal/agent"
	"other-rpc/internal/agent/embedding"
	"other-rpc/internal/agent/vector"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
)

type ServiceContext struct {
	Config   config.Config
	Embedder embeddings.Embedder
	Chroma   *vector.Chroma
	Agent    *agent.Agent
}

func NewServiceContext(c config.Config) *ServiceContext {
	embedder, err := embedding.NewEmbedder(c.KnowledgeBase.Embedding)
	if err != nil {
		fmt.Println("初始化 Embedding 模块失败:", err)
	}

	ch, err := vector.NewChroma(context.Background(), c.KnowledgeBase)
	if err != nil {
		fmt.Println("Chroma 初始化失败:", err)
	}

	retriever, err := ch.Load(context.Background(), embedder)
	if err != nil {
		fmt.Println("vector 初始化失败:", err)
	}

	kbAgent, err := agent.NewAgent(c.KnowledgeBase, retriever)
	if err != nil {
		fmt.Println("Agent 构建失败:", err)
	}

	return &ServiceContext{
		Config:   c,
		Embedder: embedder,
		Chroma:   ch,
		Agent:    kbAgent,
	}
}
