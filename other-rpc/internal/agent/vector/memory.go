package vector

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// MemoryVector 内存向量数据库实现
type MemoryVector struct {
	mu       sync.RWMutex
	embedder *embeddings.EmbedderImpl
	vectors  map[string]*vectorDoc // ID -> 向量文档
}

// vectorDoc 内部存储结构
type vectorDoc struct {
	ID      string
	Vector  []float64
	Content string
	Meta    map[string]any
}

// 确保 MemoryVector 实现了 VectorStore 接口
var _ vectorstores.VectorStore = (*MemoryVector)(nil)

// NewMemoryVector 创建内存向量数据库实例
func NewMemoryVector(embedder *embeddings.EmbedderImpl) *MemoryVector {
	return &MemoryVector{
		embedder: embedder,
		vectors:  make(map[string]*vectorDoc),
	}
}

// AddDocuments 添加文档到内存向量库
func (m *MemoryVector) AddDocuments(ctx context.Context, docs []schema.Document, opts ...vectorstores.Option) ([]string, error) {
	if len(docs) == 0 {
		return []string{}, nil
	}

	// 提取文本内容
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.PageContent
	}

	// 批量生成向量
	vectors, err := m.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("生成向量失败: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	ids := make([]string, len(docs))
	for i, doc := range docs {
		id := uuid.New().String()
		ids[i] = id

		m.vectors[id] = &vectorDoc{
			ID:      id,
			Vector:  vectors[i],
			Content: doc.PageContent,
			Meta:    doc.Metadata,
		}
	}

	return ids, nil
}

// SimilaritySearch 相似度搜索
func (m *MemoryVector) SimilaritySearch(ctx context.Context, query string, numDocuments int, opts ...vectorstores.Option) ([]schema.Document, error) {
	// 处理选项
	options := &vectorstores.Options{}
	for _, opt := range opts {
		opt(options)
	}

	// 默认返回数量
	if numDocuments <= 0 {
		numDocuments = 4
	}

	// 1. 查询向量化
	queryVec, err := m.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询向量化失败: %w", err)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.vectors) == 0 {
		return []schema.Document{}, nil
	}

	// 2. 计算所有向量的余弦相似度
	type scoredDoc struct {
		doc   *vectorDoc
		score float64
	}
	results := make([]scoredDoc, 0, len(m.vectors))

	for _, vdoc := range m.vectors {
		sim := cosineSimilarity(queryVec, vdoc.Vector)
		results = append(results, scoredDoc{doc: vdoc, score: sim})
	}

	// 3. 按相似度降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// 4. 转换为 schema.Document 并应用过滤
	docs := make([]schema.Document, 0, numDocuments)
	for i := 0; i < len(results) && len(docs) < numDocuments; i++ {
		res := results[i]

		// 应用相似度阈值过滤
		if options.ScoreThreshold > 0 && res.score < float64(options.ScoreThreshold) {
			continue
		}

		// 应用元数据过滤
		if options.Filters != nil && !matchFilters(res.doc.Meta, options.Filters) {
			continue
		}

		docs = append(docs, schema.Document{
			PageContent: res.doc.Content,
			Metadata:    res.doc.Meta,
			Score:       float32(res.score),
		})
	}

	return docs, nil
}

// ========== 辅助函数 ==========

// cosineSimilarity 计算两个向量的余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// matchFilters 检查元数据是否匹配过滤条件
// 支持的过滤格式：{"field": "value"} 或 {"field": 123}
func matchFilters(meta map[string]any, filters map[string]any) bool {
	if filters == nil || len(filters) == 0 {
		return true
	}

	for key, expected := range filters {
		actual, ok := meta[key]
		if !ok {
			return false
		}

		// 简单值比较（可扩展更多类型）
		if fmt.Sprint(actual) != fmt.Sprint(expected) {
			return false
		}
	}
	return true
}
