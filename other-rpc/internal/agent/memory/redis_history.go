// Package memory 提供 Agent 短期对话记忆的 Redis 实现。
//
// redis_history.go 实现 langchaingo 的 schema.ChatMessageHistory 接口，
// 将每个 session 的 Q&A 历史序列化为 JSON 存入 Redis，供 ConversationalRetrievalQA
// 在多轮对话时读取上下文、改写追问并检索。
package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// redisKeyPrefix 是所有会话历史 key 的公共前缀，完整 key 形如：
// agent:session:{session_id}:history
const redisKeyPrefix = "agent:session:"

// storedMessage 是写入 Redis 的 JSON 结构，只保留角色与文本内容。
// 格式：
//
//	{
//	 "messages":[
//	   {"role":"user","content":"你好"},
//	   {"role":"assistant","content":"你好，我是AI助手"},
//	   {"role":"user","content":"帮我写一个go-zero项目"}
//	 ]
//	}
type storedMessage struct {
	Role    string `json:"role"`    // human / ai / system
	Content string `json:"content"` // 消息正文
}

// RedisChatMessageHistory 按 session_id 读写 Redis 中的对话历史。

type RedisChatMessageHistory struct {
	client *redis.Redis

	// sessionID 是当前会话的唯一标识，在chat的时候构建对应的memory
	sessionID string

	// 短期记忆的记忆时间
	ttl int
}

// 编译期断言：确保实现了 langchaingo 要求的 ChatMessageHistory 接口。
var _ schema.ChatMessageHistory = (*RedisChatMessageHistory)(nil)

// NewRedisChatMessageHistory 创建绑定到指定会话的 Redis 历史存储。
func NewRedisChatMessageHistory(client *redis.Redis, sessionID string, ttlSeconds int) *RedisChatMessageHistory {
	return &RedisChatMessageHistory{
		client:    client,
		sessionID: sessionID,
		ttl:       ttlSeconds,
	}
}

// redisKey 返回当前会话在 Redis 中的 key。
func (h *RedisChatMessageHistory) redisKey() string {
	return redisKeyPrefix + h.sessionID + ":history"
}

// Messages 从 Redis 读取全部历史并还原为 langchaingo 消息列表。
func (h *RedisChatMessageHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	// 1.查找对应的 session_id 对应的 短期记忆（为空也无所谓）
	raw, err := h.client.GetCtx(ctx, h.redisKey())
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []llms.ChatMessage{}, nil
	}

	var stored []storedMessage
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return nil, fmt.Errorf("decode session history: %w", err)
	}

	msgs := make([]llms.ChatMessage, 0, len(stored))
	for _, item := range stored {
		msgs = append(msgs, fromStored(item))
	}
	return msgs, nil
}

// AddAIMessage 追加一条助手回复（RAG 问答结束后由链自动调用）。
func (h *RedisChatMessageHistory) AddAIMessage(ctx context.Context, text string) error {
	return h.addMessage(ctx, llms.AIChatMessage{Content: text})
}

// AddUserMessage 追加一条用户提问（RAG 问答开始时由链自动调用）。
func (h *RedisChatMessageHistory) AddUserMessage(ctx context.Context, text string) error {
	return h.addMessage(ctx, llms.HumanChatMessage{Content: text})
}

// AddMessage 追加任意类型的聊天消息。
func (h *RedisChatMessageHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	return h.addMessage(ctx, message)
}

// Clear 删除当前会话的全部历史（langchaingo 链在重置记忆时可能调用）。
func (h *RedisChatMessageHistory) Clear(ctx context.Context) error {
	_, err := h.client.DelCtx(ctx, h.redisKey())
	return err
}

// SetMessages 用新列表整体替换历史。
// ConversationWindowBuffer 在超出 WindowTurns 后会裁剪消息并调用此方法写回。
func (h *RedisChatMessageHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	stored := make([]storedMessage, 0, len(messages))
	for _, msg := range messages {
		stored = append(stored, toStored(msg))
	}
	return h.save(ctx, stored)
}

// addMessage 读-改-写：先拉取现有历史，追加一条，再整包写回 Redis。
func (h *RedisChatMessageHistory) addMessage(ctx context.Context, message llms.ChatMessage) error {
	msgs, err := h.Messages(ctx)
	if err != nil {
		return err
	}
	msgs = append(msgs, message)

	stored := make([]storedMessage, 0, len(msgs))
	for _, msg := range msgs {
		stored = append(stored, toStored(msg))
	}
	return h.save(ctx, stored)
}

// save 将历史序列化为 JSON 写入 Redis；ttl > 0 时使用 SETEX 并刷新过期时间。
func (h *RedisChatMessageHistory) save(ctx context.Context, stored []storedMessage) error {
	raw, err := json.Marshal(stored)
	if err != nil {
		return fmt.Errorf("encode session history: %w", err)
	}

	if h.ttl > 0 {
		return h.client.SetexCtx(ctx, h.redisKey(), string(raw), h.ttl)
	}
	return h.client.SetCtx(ctx, h.redisKey(), string(raw))
}

// toStored 将 langchaingo ChatMessage 转为可 JSON 序列化的结构。
func toStored(msg llms.ChatMessage) storedMessage {
	return storedMessage{
		Role:    string(msg.GetType()),
		Content: msg.GetContent(),
	}
}

// fromStored 将 Redis 中的 JSON 记录还原为对应的 ChatMessage 具体类型。
func fromStored(item storedMessage) llms.ChatMessage {
	switch llms.ChatMessageType(item.Role) {
	case llms.ChatMessageTypeAI:
		return llms.AIChatMessage{Content: item.Content}
	case llms.ChatMessageTypeSystem:
		return llms.SystemChatMessage{Content: item.Content}
	default:
		return llms.HumanChatMessage{Content: item.Content}
	}
}
