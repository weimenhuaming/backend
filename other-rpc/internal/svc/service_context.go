package svc

import (
	"log"
	"other-rpc/internal/agent"
	"other-rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Agent  *agent.Agent
}

func NewServiceContext(c config.Config) *ServiceContext {
	kbAgent, err := agent.NewAgent(c.KnowledgeBase)
	if err != nil {
		log.Fatalf("初始化知识库 Agent 失败: %v", err)
	}
	docCount, chunkCount, topK := kbAgent.Stats()
	logx.Infof("知识库 Agent 已就绪，Chroma: %s, collection: %s, 文档数: %d, 切片数: %d, TopK: %d",
		c.KnowledgeBase.Chroma.URL, c.KnowledgeBase.Chroma.Collection, docCount, chunkCount, topK)

	return &ServiceContext{
		Config: c,
		Agent:  kbAgent,
	}
}
