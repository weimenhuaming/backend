package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
)

// QA 封装 RAG 检索问答链。
type QA struct {
	chain chains.RetrievalQA
}

// NewQA 使用 LLM + Retriever 构建检索问答链。
func NewQA(model llms.Model, retriever vectorstores.Retriever) *QA {
	return &QA{
		chain: chains.NewRetrievalQAFromLLM(model, retriever),
	}
}

// Ask 对外提供统一问答入口。
func (q *QA) Ask(ctx context.Context, question string) (string, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return "", fmt.Errorf("问题不能为空")
	}
	return chains.Run(ctx, q.chain, question)
}
