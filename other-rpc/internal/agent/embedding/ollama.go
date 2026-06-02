package embedding

import (
	"errors"
	"fmt"
	"strings"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

// NewEmbedder 按配置初始化 Embedding 模块。
func NewEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))

	if provider != "" && provider != "ollama" {
		return nil, fmt.Errorf("不支持的 Embedding Provider: %s", cfg.Provider)
	}

	model := cfg.Model
	if model == "" {
		return nil, errors.New("embedding model is not find")
	}

	opts := []ollama.Option{ollama.WithModel(model)}
	if cfg.BaseURL != "" {
		opts = append(opts, ollama.WithServerURL(cfg.BaseURL))
	}

	embedLLM, err := ollama.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("创建 Ollama embedder 失败: %w", err)
	}

	return embeddings.NewEmbedder(embedLLM)
}
