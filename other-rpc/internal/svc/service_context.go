package svc

import (
	"context"
	"fmt"
	"log"
	"other-rpc/internal/agent"
	"other-rpc/internal/agent/embedding"
	"other-rpc/internal/agent/vector"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
		log.Fatal("初始化 Embedding 模块失败:", err)
	}
	fmt.Println("Embedding model 连接成功")

	ch, err := vector.NewChroma(context.Background(), c.KnowledgeBase)
	if err != nil {
		log.Fatal("Chroma 初始化失败:", err)
	}
	fmt.Println("Chroma 连接成功")

	retriever, err := ch.Load(context.Background(), embedder)
	if err != nil {
		log.Fatal("vector 初始化失败:", err)
	}

	cache := redis.MustNewRedis(c.Cache[0])
	fmt.Println("Redis 连接成功")

	kbAgent, err := agent.NewAgent(c.KnowledgeBase, retriever, cache)
	if err != nil {
		log.Fatal("Agent 构建失败:", err)
	}
	fmt.Println("Agent 构建成功")

	return &ServiceContext{
		Config:   c,
		Embedder: embedder,
		Chroma:   ch,
		Agent:    kbAgent,
	}
}
