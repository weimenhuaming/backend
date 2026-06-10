package rag

import (
	"context"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/vectorstores"
)

// QA 封装 RAG 检索问答链。
type QA struct {
	chain chains.RetrievalQA
}

// NewQA 使用 LLM + Retriever 构建带系统提示词的检索问答链。
func NewQA(model llms.Model, retriever vectorstores.Retriever) *QA {
	prompt := prompts.NewPromptTemplate(qaTemplate, []string{"context", "question"})
	llmChain := chains.NewLLMChain(model, prompt)
	combineChain := chains.NewStuffDocuments(llmChain)

	return &QA{
		chain: chains.NewRetrievalQA(combineChain, retriever),
	}
}

// Ask 对外提供统一问答入口。
func (q *QA) Ask(ctx context.Context, question string) (string, error) {
	return chains.Run(ctx, q.chain, question)
}

// AskStream 流式问答，每收到 LLM token 即回调 send。
func (q *QA) AskStream(ctx context.Context, question string, send func(chunk string) error) error {
	_, err := chains.Run(ctx, q.chain, question, chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		if len(chunk) == 0 {
			return nil
		}
		return send(string(chunk))
	}))
	return err
}
