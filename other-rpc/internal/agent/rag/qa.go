package rag

import (
	"context"
	"other-rpc/internal/agent/memory"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// QA 封装带短期记忆的 RAG 检索问答链。
type QA struct {
	model     llms.Model
	retriever vectorstores.Retriever
	sessions  *memory.Store
}

// NewQA 使用 LLM + Retriever 构建带系统提示词与短期会话记忆的问答链。
func NewQA(model llms.Model, retriever vectorstores.Retriever, sessions *memory.Store) *QA {
	return &QA{
		model:     model,
		retriever: retriever,
		sessions:  sessions,
	}
}

func (q *QA) buildChain(sessionMemory schema.Memory) chains.ConversationalRetrievalQA {
	prompt := prompts.NewPromptTemplate(qaTemplate, []string{"context", "question"})
	llmChain := chains.NewLLMChain(q.model, prompt)
	combineChain := chains.NewStuffDocuments(llmChain)

	condensePrompt := prompts.NewPromptTemplate(condenseQuestionTemplate, []string{"chat_history", "question"})
	condenseChain := chains.NewLLMChain(q.model, condensePrompt)

	return chains.NewConversationalRetrievalQA(
		combineChain,
		condenseChain,
		q.retriever,
		sessionMemory,
	)
}

// Ask 对外提供统一问答入口，sessionID 相同则带上短期对话上下文。
func (q *QA) Ask(ctx context.Context, sessionID, question string) (string, error) {
	chain := q.buildChain(q.sessions.Memory(sessionID))
	return chains.Run(ctx, chain, question)
}

// AskStream 流式问答，每收到 LLM token 即回调 send。
func (q *QA) AskStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	chain := q.buildChain(q.sessions.Memory(sessionID))
	_, err := chains.Run(ctx, chain, question, chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		if len(chunk) == 0 {
			return nil
		}
		return send(string(chunk))
	}))
	return err
}
