// 调用方（rag.QA）每次问答前通过 Store.Memory(sessionID) 取到该会话的记忆对象；
// langchaingo 的 ConversationalRetrievalQA 会在问答过程中自动读写该对象。
package memory

import (
	"other-rpc/internal/config"
	"strings"
	"sync"
	"time"

	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

type Store struct {
	mu          sync.RWMutex
	sessions    map[string]*entry
	windowTurns int
	ttl         time.Duration
}

type entry struct {
	mem      schema.Memory // langchaingo 的对话窗口缓冲，问答链会直接读写
	lastSeen time.Time     // 最近一次 Memory(sessionID) 被调用的时间
}

func NewStore(cfg config.MemoryConf) *Store {
	windowTurns := cfg.WindowTurns
	if windowTurns <= 0 {
		windowTurns = 3
	}
	ttl := time.Duration(cfg.SessionTTL) * time.Second
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	return &Store{
		sessions:    make(map[string]*entry),
		windowTurns: windowTurns,
		ttl:         ttl,
	}
}

// Memory 返回指定会话的记忆对象，供 RAG 链挂载使用。
func (s *Store) Memory(sessionID string) schema.Memory {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		// 无 session：不跨请求保留历史
		return s.newBuffer()
	}

	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.evictExpiredLocked(now)

	if e, ok := s.sessions[sessionID]; ok {
		e.lastSeen = now
		return e.mem
	}

	// 新会话：创建空缓冲，后续问答会自动追加 Q&A
	buf := s.newBuffer()
	s.sessions[sessionID] = &entry{mem: buf, lastSeen: now}
	return buf
}

// newBuffer 创建滑动窗口对话缓冲。
// InputKey/OutputKey 须与 rag 链传入的字段名一致（question → text）。
func (s *Store) newBuffer() schema.Memory {
	return memory.NewConversationWindowBuffer(
		s.windowTurns,
		memory.WithInputKey("question"),
		memory.WithOutputKey("text"),
	)
}

// evictExpiredLocked 删除长时间未访问的会话。调用方需已持有 s.mu。
func (s *Store) evictExpiredLocked(now time.Time) {
	for id, e := range s.sessions {
		if now.Sub(e.lastSeen) > s.ttl {
			delete(s.sessions, id)
		}
	}
}
