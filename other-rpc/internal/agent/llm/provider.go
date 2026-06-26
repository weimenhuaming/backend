package llm

import (
	"fmt"
	"strings"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// NewChatModel 按配置初始化对话模型。
func NewChatModel(cfg config.LLMConf) (llms.Model, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))

	// 根据对应的名字选择对应llm模型的构造器
	switch provider {
	case "openai":
		return newOpenAIChatModel(cfg)
	default:
		return nil, fmt.Errorf("不支持的 Chat LLM Provider: %s", cfg.Provider)
	}
}

// newOpenAIChatModel 目前只支持openai接口的
func newOpenAIChatModel(cfg config.LLMConf) (llms.Model, error) {
	opts := []openai.Option{}
	if cfg.APIKey != "" {
		opts = append(opts, openai.WithToken(cfg.APIKey))
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
	}
	if cfg.Model != "" {
		opts = append(opts, openai.WithModel(cfg.Model))
	}
	chatLLM, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}
	return chatLLM, nil
}
