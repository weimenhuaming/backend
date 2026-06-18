package embedding

import (
	"errors"
	"fmt"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func newOllamaEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	if cfg.Model == "" {
		return nil, errors.New("embedding model 未配置")
	}
	if cfg.BaseURL == "" {
		return nil, errors.New("embedding BaseURL 未配置")
	}

	opts := []ollama.Option{
		ollama.WithModel(cfg.Model),
		ollama.WithServerURL(cfg.BaseURL),
	}

	embedLLM, err := ollama.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("创建 Ollama embedder 失败: %w", err)
	}

	return embeddings.NewEmbedder(embedLLM)
}
