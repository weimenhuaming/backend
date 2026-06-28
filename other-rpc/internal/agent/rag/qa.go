package rag

import (
	"context"
	"fmt"
	"strings"

	"other-rpc/internal/agent/memory"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// QA 封装带短期记忆的 RAG 检索问答。
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

// Ask 同步问答。
func (q *QA) Ask(ctx context.Context, sessionID, question string) (string, error) {
	answer, err := q.run(ctx, sessionID, question, nil)
	if err != nil {
		return "", err
	}
	return answer, nil
}

// AskStream 流式问答，每收到 LLM token 即回调 send。
func (q *QA) AskStream(ctx context.Context, sessionID, question string, send func(chunk string) error) error {
	_, err := q.run(ctx, sessionID, question, send)
	return err
}

func (q *QA) run(
	ctx context.Context,
	sessionID, question string,
	send func(chunk string) error,
) (string, error) {
	history, err := q.sessions.LoadHistoryText(ctx, sessionID)
	if err != nil {
		return "", err
	}

	lastUserQuestion, err := q.sessions.LastHumanQuestion(ctx, sessionID)
	if err != nil {
		return "", err
	}

	if answer, ok := answerPreviousQuestion(question, lastUserQuestion); ok {
		if err := q.emitAnswer(send, answer); err != nil {
			return "", err
		}
		if err := q.sessions.AppendTurn(ctx, sessionID, question, answer); err != nil {
			return "", err
		}
		return answer, nil
	}

	retrievalQuestion, err := q.rewriteForRetrieval(ctx, history, question)
	if err != nil {
		return "", err
	}

	docs, err := q.retriever.GetRelevantDocuments(ctx, retrievalQuestion)
	if err != nil {
		return "", err
	}

	prompt, err := formatQAPrompt(history, lastUserQuestion, joinDocuments(docs), question)
	if err != nil {
		return "", err
	}

	opts := make([]llms.CallOption, 0, 1)
	if send != nil {
		opts = append(opts, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) == 0 {
				return nil
			}
			return send(string(chunk))
		}))
	}

	answer, err := llms.GenerateFromSinglePrompt(ctx, q.model, prompt, opts...)
	if err != nil {
		return "", err
	}

	if err := q.sessions.AppendTurn(ctx, sessionID, question, answer); err != nil {
		return "", err
	}

	return answer, nil
}

// answerPreviousQuestion 对「上一问题是什么」类追问，直接从 Redis 历史取答案，避免模型混淆当前问题。
func answerPreviousQuestion(question, lastUserQuestion string) (string, bool) {
	if !isAskPreviousQuestion(question) {
		return "", false
	}
	if lastUserQuestion == "" {
		return "这是本会话的第一个问题，目前还没有上一个问题。", true
	}
	return fmt.Sprintf("您的上一个问题是：「%s」。", lastUserQuestion), true
}

func isAskPreviousQuestion(question string) bool {
	q := strings.TrimSpace(question)
	for _, kw := range []string{"上一个问题", "上一题", "刚才问", "之前问", "刚才的问题", "之前的问题"} {
		if strings.Contains(q, kw) {
			return true
		}
	}
	return false
}

func (q *QA) emitAnswer(send func(chunk string) error, answer string) error {
	if send == nil {
		return nil
	}
	return send(answer)
}

func (q *QA) rewriteForRetrieval(ctx context.Context, history, question string) (string, error) {
	if strings.TrimSpace(history) == "" {
		return question, nil
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

func joinDocuments(docs []schema.Document) string {
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

var (
	qaPrompt       = prompts.NewPromptTemplate(qaTemplate, []string{"chat_history", "last_user_question", "context", "question"})
	condensePrompt = prompts.NewPromptTemplate(condenseQuestionTemplate, []string{"chat_history", "question"})
)

func formatQAPrompt(history, lastUserQuestion, context, question string) (string, error) {
	if lastUserQuestion == "" {
		lastUserQuestion = "（无，这是本会话第一个问题）"
	}

	promptValue, err := qaPrompt.FormatPrompt(map[string]any{
		"chat_history":       history,
		"last_user_question": lastUserQuestion,
		"context":            context,
		"question":           question,
	})
	if err != nil {
		return "", err
	}
	return promptValue.String(), nil
}
