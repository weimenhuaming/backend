package kb

import (
	"context"
	"fmt"

	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

// Index 表示知识库构建后的可检索结果。
type Index struct {
	Store      *MemoryStore
	Retriever  vectorstores.Retriever
	TopK       int
	DocCount   int
	ChunkCount int
}

// Build 完成知识库加载、切分、向量化与检索器构建。
func Build(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (*Index, error) {
	store := NewMemoryStore(embedder)

	docs, err := LoadDocumentsFromDir(cfg.DataPath)
	if err != nil {
		return nil, fmt.Errorf("加载知识库目录失败: %w", err)
	}

	chunkCount, err := indexDocuments(ctx, store, cfg, docs)
	if err != nil {
		return nil, err
	}

	topK := cfg.TopK
	if topK <= 0 {
		topK = 4
	}

	return &Index{
		Store:      store,
		Retriever:  vectorstores.ToRetriever(store, topK),
		TopK:       topK,
		DocCount:   len(docs),
		ChunkCount: chunkCount,
	}, nil
}

func indexDocuments(ctx context.Context, store *MemoryStore, cfg config.KnowledgeBaseConf, docs []schema.Document) (int, error) {
	if len(docs) == 0 {
		return 0, nil
	}

	chunkSize := cfg.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 800
	}
	chunkOverlap := cfg.ChunkOverlap
	if chunkOverlap <= 0 {
		chunkOverlap = 100
	}
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(chunkSize),
		textsplitter.WithChunkOverlap(chunkOverlap),
	)
	chunks, err := textsplitter.SplitDocuments(splitter, docs)
	if err != nil {
		return 0, fmt.Errorf("切分文档失败: %w", err)
	}
	if len(chunks) == 0 {
		return 0, nil
	}
	if _, err = store.AddDocuments(ctx, chunks); err != nil {
		return 0, err
	}
	return len(chunks), nil
}
