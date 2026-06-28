package memory

import (
	"context"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Store struct {
	redis       *redis.Redis
	windowTurns int
	ttlSeconds  int
}

func NewStore(cfg config.MemoryConf, client *redis.Redis) *Store {
	return &Store{
		redis:       client,
		windowTurns: cfg.WindowTurns,
		ttlSeconds:  cfg.SessionTTL,
	}
}

// LoadHistoryText 读取当前会话已有问答，格式化为 Human/AI 文本供 Prompt 使用。
// 不包含本轮尚未写入的问题。
func (s *Store) LoadHistoryText(ctx context.Context, sessionID string) (string, error) {
	msgs, err := s.loadMessages(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if len(msgs) == 0 {
		return "", nil
	}
	return llms.GetBufferString(msgs, "Human", "AI")
}

// LastHumanQuestion 返回对话历史中最后一条用户消息（即上一轮用户问题）。
// 本轮 question 尚未写入 Redis，因此不会与当前问题混淆。
func (s *Store) LastHumanQuestion(ctx context.Context, sessionID string) (string, error) {
	msgs, err := s.loadMessages(ctx, sessionID)
	if err != nil {
		return "", err
	}
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].GetType() == llms.ChatMessageTypeHuman {
			return msgs[i].GetContent(), nil
		}
	}
	return "", nil
}

func (s *Store) loadMessages(ctx context.Context, sessionID string) ([]llms.ChatMessage, error) {
	if sessionID == "" {
		return nil, nil
	}
	history := NewRedisChatMessageHistory(s.redis, sessionID, s.ttlSeconds)
	return history.Messages(ctx)
}

// AppendTurn 将本轮问答追加到 Redis，并按窗口大小裁剪。
func (s *Store) AppendTurn(ctx context.Context, sessionID, question, answer string) error {
	if sessionID == "" {
		return nil
	}

	history := NewRedisChatMessageHistory(s.redis, sessionID, s.ttlSeconds)
	msgs, err := history.Messages(ctx)
	if err != nil {
		return err
	}

	msgs = append(msgs,
		llms.HumanChatMessage{Content: question},
		llms.AIChatMessage{Content: answer},
	)
	msgs = trimMessages(msgs, s.windowTurns)
	return history.SetMessages(ctx, msgs)
}

func trimMessages(msgs []llms.ChatMessage, windowTurns int) []llms.ChatMessage {
	if windowTurns <= 0 {
		windowTurns = 5
	}
	maxMessages := windowTurns * 2
	if len(msgs) <= maxMessages {
		return msgs
	}
	return msgs[len(msgs)-maxMessages:]
}
