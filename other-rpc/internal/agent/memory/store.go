package memory

import (
	"strings"
	"sync"
	"time"

	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

// Store 按 session_id 保存短期对话记忆（进程内，重启后丢失）。
type Store struct {
	mu          sync.RWMutex
	sessions    map[string]*entry
	windowTurns int
	ttl         time.Duration
}

type entry struct {
	mem      schema.Memory
	lastSeen time.Time
}

// NewStore 创建会话记忆存储。windowTurns 为保留的最近轮数，ttl 为会话空闲过期时间。
func NewStore(windowTurns int, ttl time.Duration) *Store {
	if windowTurns <= 0 {
		windowTurns = 5
	}
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return &Store{
		sessions:    make(map[string]*entry),
		windowTurns: windowTurns,
		ttl:         ttl,
	}
}

// Memory 返回指定会话的记忆；sessionID 为空时使用不落库的临时记忆。
func (s *Store) Memory(sessionID string) schema.Memory {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
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

	buf := s.newBuffer()
	s.sessions[sessionID] = &entry{mem: buf, lastSeen: now}
	return buf
}

func (s *Store) newBuffer() schema.Memory {
	return memory.NewConversationWindowBuffer(
		s.windowTurns,
		memory.WithInputKey("question"),
		memory.WithOutputKey("text"),
	)
}

func (s *Store) evictExpiredLocked(now time.Time) {
	for id, e := range s.sessions {
		if now.Sub(e.lastSeen) > s.ttl {
			delete(s.sessions, id)
		}
	}
}
