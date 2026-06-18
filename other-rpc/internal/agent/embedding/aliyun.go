package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"other-rpc/internal/config"
)

const (
	//  单次请求最多 20 条内容。
	aliyunMaxBatchSize = 20
)

type aliyunEmbedder struct {
	apiKey    string
	baseURL   string
	model     string
	dimension int
	client    *http.Client
}

type aliyunEmbeddingRequest struct {
	Model      string                 `json:"model"`
	Input      aliyunEmbeddingInput   `json:"input"`
	Parameters *aliyunEmbeddingParams `json:"parameters,omitempty"`
}

type aliyunEmbeddingInput struct {
	Contents []map[string]string `json:"contents"`
}

type aliyunEmbeddingParams struct {
	Dimension int `json:"dimension,omitempty"`
}

type aliyunEmbeddingResponse struct {
	Output struct {
		Embeddings []struct {
			Index     int       `json:"index"`
			Embedding []float64 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newAliyunEmbedder(cfg config.EmbeddingConf) (*aliyunEmbedder, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("embedding APIKey 未配置")
	}
	if cfg.Model == "" {
		return nil, errors.New("embedding model 未配置")
	}
	if cfg.BaseURL == "" {
		return nil, errors.New("embedding BaseURL 未配置")
	}
	if cfg.Dimension <= 0 {
		return nil, errors.New("embedding dimension 未配置")
	}

	return &aliyunEmbedder{
		apiKey:    cfg.APIKey,
		baseURL:   cfg.BaseURL,
		model:     cfg.Model,
		dimension: cfg.Dimension,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (e *aliyunEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	vectors := make([][]float32, len(texts))
	for start := 0; start < len(texts); start += aliyunMaxBatchSize {
		end := start + aliyunMaxBatchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch, err := e.embedTexts(ctx, texts[start:end])
		if err != nil {
			return nil, err
		}
		copy(vectors[start:end], batch)
	}
	return vectors, nil
}

func (e *aliyunEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	vectors, err := e.embedTexts(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vectors[0], nil
}

func (e *aliyunEmbedder) embedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	contents := make([]map[string]string, len(texts))
	for i, text := range texts {
		contents[i] = map[string]string{"text": text}
	}

	reqBody := aliyunEmbeddingRequest{
		Model: e.model,
		Input: aliyunEmbeddingInput{Contents: contents},
		Parameters: &aliyunEmbeddingParams{
			Dimension: e.dimension,
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化阿里云 Embedding 请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("创建阿里云 Embedding 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用阿里云 Embedding 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取阿里云 Embedding 响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("阿里云 Embedding 返回 HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result aliyunEmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析阿里云 Embedding 响应失败: %w", err)
	}
	if result.Code != "" {
		return nil, fmt.Errorf("阿里云 Embedding 错误 [%s]: %s", result.Code, result.Message)
	}
	if len(result.Output.Embeddings) != len(texts) {
		return nil, fmt.Errorf("阿里云 Embedding 返回向量数量不匹配: 期望 %d, 实际 %d", len(texts), len(result.Output.Embeddings))
	}

	sort.Slice(result.Output.Embeddings, func(i, j int) bool {
		return result.Output.Embeddings[i].Index < result.Output.Embeddings[j].Index
	})

	vectors := make([][]float32, len(texts))
	for i, item := range result.Output.Embeddings {
		if item.Index != i {
			return nil, fmt.Errorf("阿里云 Embedding 返回索引不连续: 期望 %d, 实际 %d", i, item.Index)
		}
		vectors[i] = make([]float32, len(item.Embedding))
		for j, v := range item.Embedding {
			vectors[i][j] = float32(v)
		}
	}
	return vectors, nil
}
