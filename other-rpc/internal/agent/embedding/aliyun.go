// 实现 langchaingo 的 embeddings.Embedder 接口，供 Chroma 向量库统一调用。

//	curl -X POST 'https://dashscope.aliyuncs.com/api/v1/services/embeddings/multimodal-embedding/multimodal-embedding' \
//	  -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
//	  -H "Content-Type: application/json" \
//	  -d '{
//	    "model": "qwen3-vl-embedding",
//	    "input": {"contents": [{"text": "示例文本"}]},
//	    "parameters": {"dimension": 1024}
//	  }'
//
// ├─ 建库时：EmbedDocuments() 批量把文档切块转成向量写入 Chroma
// └─ 检索时：EmbedQuery() 把用户问题转成向量，在 Chroma 中做相似度搜索
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
	aliyunMaxBatchSize = 20
)

type aliyunEmbedder struct {
	apiKey    string       // 对应 curl 的 Authorization: Bearer {apiKey}
	baseURL   string       // 对应 curl 的 POST 地址
	model     string       // 请求体 model 字段，如 qwen3-vl-embedding
	dimension int          // 请求体 parameters.dimension，须与 Chroma 集合维度一致
	client    *http.Client // 底层 HTTP 客户端，替代 curl 命令
}

// aliyunEmbeddingRequest 对应 curl --data 中的 JSON body。
type aliyunEmbeddingRequest struct {
	Model      string                 `json:"model"`
	Input      aliyunEmbeddingInput   `json:"input"`
	Parameters *aliyunEmbeddingParams `json:"parameters,omitempty"`
}

// aliyunEmbeddingInput 对应 input 字段。
// 每条 content 是一个 map，纯文本场景只需 {"text": "..."}；
// 多模态场景还可传 {"image": "url"}、{"video": "url"} 等（当前 RAG 仅用 text）。
type aliyunEmbeddingInput struct {
	Contents []map[string]string `json:"contents"`
}

// aliyunEmbeddingParams 对应 parameters 字段。dimension 指定输出向量维度
type aliyunEmbeddingParams struct {
	Dimension int `json:"dimension,omitempty"`
}

// aliyunEmbeddingResponse 对应 API 返回的 JSON body。
// 成功时 output.embeddings 包含向量数组；失败时 code/message 非空。
type aliyunEmbeddingResponse struct {
	Output struct {
		Embeddings []struct {
			Index     int       `json:"index"`     // 与请求 contents 的下标对应
			Embedding []float64 `json:"embedding"` // 浮点向量，API 返回 float64，此处转为 float32 供 Chroma 使用
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

// EmbedDocuments 批量向量化文档，实现 embeddings.Embedder 接口。
// 使用场景：建库（Build）时，Chroma 把知识库文档切块后调用此方法，
// 将每段文本转为 float32 向量再写入向量数据库。
// 见 vector/collection.go → AddDocuments()。
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

// EmbedQuery 向量化单条查询，实现 embeddings.Embedder 接口。
// 使用场景：RAG 检索时，Chroma 把用户问题转为向量，
// 再在库中做余弦相似度搜索，找出最相关的文档块。
// 见 vector/collection.go → SimilaritySearch()。
func (e *aliyunEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	vectors, err := e.embedTexts(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vectors[0], nil
}

// embedTexts 核心 HTTP 调用，等价于执行一次 curl POST 请求。
// 流程：组装 JSON body → POST → 读响应 → 解析向量 → 按 index 排序校验。
func (e *aliyunEmbedder) embedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	// 1. 组装 input.contents：每条文本对应一个 {"text": "..."}
	contents := make([]map[string]string, len(texts))
	for i, text := range texts {
		contents[i] = map[string]string{"text": text}
	}

	// 2. 构造请求体（对应 curl --data 的 JSON）
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

	// 3. 发起 HTTP POST（对应 curl -X POST ...）
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

	// 4. 读取并解析响应
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

	// 5. 按 index 排序，确保返回顺序与输入 texts 一一对应
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
