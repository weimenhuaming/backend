package embedding

import (
	"fmt"
	"strings"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// NewEmbedder 按配置初始化 Embedding 模块。
func NewEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	switch provider {
	case "", "ollama":
		return newOllamaEmbedder(cfg)
	case "openai":
		return newOpenAIEmbedder(cfg)
	default:
		return nil, fmt.Errorf("不支持的 Embedding Provider: %s", cfg.Provider)
	}
}

func newOllamaEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	model := cfg.Model
	if model == "" {
		model = "nomic-embed-text"
	}
	opts := []ollama.Option{ollama.WithModel(model)}
	if cfg.BaseURL != "" {
		opts = append(opts, ollama.WithServerURL(cfg.BaseURL))
	}
	embedLLM, err := ollama.New(opts...)
	if err != nil {
		return nil, err
	}
	return embeddings.NewEmbedder(embedLLM)
}

func newOpenAIEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	opts := []openai.Option{}
	if cfg.APIKey != "" {
		opts = append(opts, openai.WithToken(cfg.APIKey))
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
	}
	model := cfg.Model
	if model == "" {
		model = "text-embedding-3-small"
	}
	opts = append(opts, openai.WithEmbeddingModel(model))
	embedLLM, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}
	return embeddings.NewEmbedder(embedLLM)
}
