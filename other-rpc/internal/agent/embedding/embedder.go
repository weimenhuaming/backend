package embedding

import (
	"errors"
	"fmt"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
)

// NewEmbedder 按配置初始化 Embedding 模块，所有必填项均须由外部配置显式提供。
func NewEmbedder(cfg config.EmbeddingConf) (embeddings.Embedder, error) {
	if cfg.Provider == "" {
		return nil, errors.New("embedding provider 未配置")
	}

	switch cfg.Provider {
	case "ollama":
		return newOllamaEmbedder(cfg)
	case "aliyun", "dashscope":
		return newAliyunEmbedder(cfg)
	default:
		return nil, fmt.Errorf("不支持的 Embedding Provider: %s", cfg.Provider)
	}
}
