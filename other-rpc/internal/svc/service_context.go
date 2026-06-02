package svc

import (
	"context"
	"fmt"
	"other-rpc/internal/agent"
	"other-rpc/internal/agent/embedding"
	"other-rpc/internal/agent/vector"
	"other-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Agent  *agent.Agent
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 拿到embedding构造器
	embedder, err := embedding.NewEmbedder(c.KnowledgeBase.Embedding)
	if err != nil {
		fmt.Println("初始化 Embedding 模块失败:", err)
	}

	// 2.拿到检索器
	retriever, err := vector.Load(context.Background(), c.KnowledgeBase, embedder)
	if err != nil {
		fmt.Println("vector 初始化失败:", err)
	}

	// 3.构建agent
	kbAgent, err := agent.NewAgent(c.KnowledgeBase, embedder, retriever)
	if err != nil {
		fmt.Println("Agent 构建失败:", err)
	}

	return &ServiceContext{
		Config: c,
		Agent:  kbAgent,
	}
}
