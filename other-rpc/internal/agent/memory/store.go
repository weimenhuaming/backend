// 调用方（rag.QA）每次问答前通过 Store.Memory(sessionID) 取到该会话的记忆对象；
// langchaingo 的 ConversationalRetrievalQA 会在问答过程中自动读写该对象。
package memory

import (
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
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

// Memory 返回指定会话的记忆对象，供 RAG 链挂载使用。
func (s *Store) Memory(sessionID string) schema.Memory {
	//  1. sessionID 为空 → 返回临时缓冲，不跨请求保留历史
	if sessionID == "" {
		return s.newBuffer(nil)
	}

	//  2. sessionID 非空 → 使用 Redis 缓存读写，多实例共享
	history := NewRedisChatMessageHistory(s.redis, sessionID, s.ttlSeconds)
	return s.newBuffer(history)
}

func (s *Store) newBuffer(history schema.ChatMessageHistory) schema.Memory {
	// WithInputKey / WithOutputKey 必须与 RAG 链 SaveContext 时传入的字段名一致：
	//   - 用户提问 → inputValues["question"]
	//   - 模型回答 → outputValues["text"]
	// 若 key 对不上，SaveContext 会报 ErrInvalidInputValues，历史写不进 Redis。
	opts := []memory.ConversationBufferOption{
		memory.WithInputKey("question"),
		memory.WithOutputKey("text"),
	}

	if history != nil {
		opts = append(opts, memory.WithChatHistory(history))
	}

	// 超出窗口后，Buffer 会调用 history.SetMessages 裁剪并写回 Redis。
	return memory.NewConversationWindowBuffer(s.windowTurns, opts...)
}
