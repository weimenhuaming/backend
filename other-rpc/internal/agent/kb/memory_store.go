package kb

import (
	"context"
	"math"
	"sort"
	"sync"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// MemoryStore 基于内存向量的轻量知识库，实现 langchaingo VectorStore。
type MemoryStore struct {
	embedder embeddings.Embedder
	mu       sync.RWMutex
	items    []storedDoc
}

type storedDoc struct {
	doc    schema.Document
	vector []float32
}

func NewMemoryStore(embedder embeddings.Embedder) *MemoryStore {
	return &MemoryStore{embedder: embedder}
}

func (s *MemoryStore) AddDocuments(ctx context.Context, docs []schema.Document, _ ...vectorstores.Option) ([]string, error) {
	if len(docs) == 0 {
		return nil, nil
	}
	texts := make([]string, len(docs))
	for i, d := range docs {
		texts[i] = d.PageContent
	}
	vectors, err := s.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	ids := make([]string, len(docs))
	for i, d := range docs {
		ids[i] = ""
		s.items = append(s.items, storedDoc{doc: d, vector: vectors[i]})
	}
	return ids, nil
}

func (s *MemoryStore) SimilaritySearch(ctx context.Context, query string, numDocuments int, _ ...vectorstores.Option) ([]schema.Document, error) {
	if numDocuments <= 0 {
		numDocuments = 4
	}
	qv, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		return nil, nil
	}

	type scored struct {
		doc   schema.Document
		score float32
	}
	scores := make([]scored, 0, len(s.items))
	for _, it := range s.items {
		scores = append(scores, scored{
			doc:   it.doc,
			score: cosineSimilarity(qv, it.vector),
		})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})
	if len(scores) > numDocuments {
		scores = scores[:numDocuments]
	}
	out := make([]schema.Document, len(scores))
	for i, sc := range scores {
		out[i] = sc.doc
		out[i].Score = sc.score
	}
	return out, nil
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float32
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(na))) * float32(math.Sqrt(float64(nb))))
}
