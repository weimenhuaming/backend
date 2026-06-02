package svc

import (
	"context"
	"fmt"
	"log"
	"os"
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
		fmt.Printf("初始化 Embedding 模块失败: %v", err)
		os.Exit(1)
	}

	// 2.拿到检索器
	retriever, _, _, err := vector.Load(context.Background(), c.KnowledgeBase, embedder)
	if err != nil {
		log.Fatalf("初始化检索器失败: %v", err)
	}

	kbAgent, err := agent.NewAgent(c.KnowledgeBase, embedder, retriever)
	if err != nil {
		log.Fatalf("初始化知识库 Agent 失败: %v", err)
	}

	return &ServiceContext{
		Config: c,
		Agent:  kbAgent,
	}
}
