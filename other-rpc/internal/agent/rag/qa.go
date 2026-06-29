package rag

import (
	"context"
	"other-rpc/internal/agent/memory"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// 3个模板:
var (
	simpleQAPrompt = prompts.NewPromptTemplate(simpleQATemplate, []string{"context", "question"})
	qaPrompt       = prompts.NewPromptTemplate(qaTemplate, []string{"chat_history", "context", "question"})
	condensePrompt = prompts.NewPromptTemplate(condenseQuestionTemplate, []string{"chat_history", "question"})
)

// QA 封装带短期记忆的 RAG 检索问答。
type QA struct {
	model llms.Model
	// 检索器内容
	retriever vectorstores.Retriever
	// 短期记忆
	chatHistory *memory.RedisChatMessageHistory
}

// NewQA 使用 LLM + Retriever 构建问答链；流式场景需传入绑定 session 的 chatHistory。
func NewQA(model llms.Model, retriever vectorstores.Retriever, chatHistory *memory.RedisChatMessageHistory) *QA {
	return &QA{
		model:       model,
		retriever:   retriever,
		chatHistory: chatHistory,
	}
}

// Ask 同步问答：检索知识库后直接调用 LLM，不使用短期记忆。
func (q *QA) Ask(ctx context.Context, question string) (string, error) {
	docs, err := q.retriever.GetRelevantDocuments(ctx, question)
	if err != nil {
		return "", err
	}

	promptValue, err := simpleQAPrompt.FormatPrompt(map[string]any{
		"context":  toDocuments(docs),
		"question": question,
	})
	if err != nil {
		return "", err
	}

	return llms.GenerateFromSinglePrompt(ctx, q.model, promptValue.String())
}

//======================================================================================================================

// AskStream 流式问答，每收到 LLM token 即回调 send。
func (q *QA) AskStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	q.chatHistory.SessionID = sessionID

	// 1.拿到短期记忆的内容，转为话string格式
	memoryMsg, err := q.chatHistory.Messages(ctx)
	if err != nil {
		return err
	}

	// 2.问题改写或者上下文压缩 => 处理完是短期记忆
	//retrievalQuestion, err := q.summary(ctx, memoryMsg, question)
	//if err != nil {
	//	return err
	//}

	// 3，通过检索器拿到匹配的文章
	docs, err := q.retriever.GetRelevantDocuments(ctx, question)
	if err != nil {
		return err
	}

	// 4.组装最后的提示词
	prompt, err := formatQAPrompt(memoryMsg, toDocuments(docs), question)
	if err != nil {
		return err
	}

	// 5.生成回复，并流式返回
	answer, err := llms.GenerateFromSinglePrompt(ctx, q.model, prompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				return nil
			}
			return send(string(chunk))
		}),
	)
	if err != nil {
		return err
	}

	// 6.最后键本轮回答追加到历史，需要裁剪
	memoryMsg = append(memoryMsg,
		llms.HumanChatMessage{Content: question},
		llms.AIChatMessage{Content: answer},
	)
	return q.chatHistory.SetMessages(ctx, memoryMsg)
}

// summary 对历史记忆进行压缩
func (q *QA) summary(ctx context.Context, memory []llms.ChatMessage, question string) (string, error) {
	history, err := llms.GetBufferString(memory, "Human", "AI")
	if err != nil {
		return "", err
	}

	promptValue, err := condensePrompt.FormatPrompt(map[string]any{
		"chat_history": history,
		"question":     question,
	})
	if err != nil {
		return question, err
	}

	rewritten, err := llms.GenerateFromSinglePrompt(ctx, q.model, promptValue.String())
	if err != nil {
		return question, err
	}

	rewritten = strings.TrimSpace(rewritten)
	if rewritten == "" {
		return question, nil
	}
	return rewritten, nil
}

// toDocuments 列表拼接为一个字符串，文档之间用两个换行符分隔。
func toDocuments(docs []schema.Document) string {
	if len(docs) == 0 {
		return ""
	}

	parts := make([]string, 0, len(docs))
	for _, doc := range docs {
		if text := strings.TrimSpace(doc.PageContent); text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n\n")
}

func formatQAPrompt(memory []llms.ChatMessage, context, question string) (string, error) {
	history, err := llms.GetBufferString(memory, "Human", "AI")
	if err != nil {
		return "", err
	}

	promptValue, err := qaPrompt.FormatPrompt(map[string]any{
		"chat_history": history,
		"context":      context,
		"question":     question,
	})
	if err != nil {
		return "", err
	}
	return promptValue.String(), nil
}
