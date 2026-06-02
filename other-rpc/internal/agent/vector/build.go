package vector

import (
	"context"
	"fmt"
	"other-rpc/internal/config"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

func newRetriever(store vectorstores.VectorStore, cfg config.KnowledgeBaseConf) vectorstores.Retriever {
	topK := cfg.TopK
	if topK <= 0 {
		topK = 4
	}
	return vectorstores.ToRetriever(store, topK)
}

// Load 连接 Chroma 并返回检索器（启动时使用，不会重新向量化）。
func Load(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (vectorstores.Retriever, int, int, error) {
	docCount, chunkCount, err := readCollectionStats(ctx, cfg)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, err
	}
	if chunkCount == 0 {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("Chroma collection 为空，请先调用 Build 构建知识库")
	}

	store, err := openChromaStore(ctx, cfg, embedder)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("连接 Chroma 失败: %w", err)
	}
	return newRetriever(store, cfg), docCount, chunkCount, nil
}

// BuildWithEmbedder 使用已有 Embedder 构建 Chroma 向量索引并返回检索器。
func BuildWithEmbedder(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (vectorstores.Retriever, int, int, error) {
	if err := resetChromaCollection(ctx, cfg); err != nil {
		return vectorstores.Retriever{}, 0, 0, err
	}

	store, err := openChromaStore(ctx, cfg, embedder)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("创建 Chroma collection 失败: %w", err)
	}

	docs, err := LoadDocumentsFromDir(cfg.DataPath)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("加载知识库目录失败: %w", err)
	}

	chunkCount, err := indexDocuments(ctx, store, cfg, docs)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, err
	}

	docCount := len(docs)
	if err := writeCollectionStats(ctx, cfg, docCount, chunkCount); err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("更新 Chroma collection 元数据失败: %w", err)
	}

	return newRetriever(store, cfg), docCount, chunkCount, nil
}

func indexDocuments(ctx context.Context, store vectorstores.VectorStore, cfg config.KnowledgeBaseConf, docs []schema.Document) (int, error) {
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

// BuildIndex 使用已有 Embedder 构建 Chroma 向量索引。
// 返回值: 文档数量, 切片数量, 错误
func BuildIndex(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (int, int, error) {
	// 1. 重置 collection（清空旧数据）
	if err := resetChromaCollection(ctx, cfg); err != nil {
		return 0, 0, err
	}

	// 2. 打开/创建 Chroma collection
	store, err := openChromaStore(ctx, cfg, embedder)
	if err != nil {
		return 0, 0, fmt.Errorf("创建 Chroma collection 失败: %w", err)
	}

	// 3. 加载文档
	docs, err := LoadDocumentsFromDir(cfg.DataPath)
	if err != nil {
		return 0, 0, fmt.Errorf("加载知识库目录失败: %w", err)
	}

	// 4. 索引文档（分块、向量化、存储）
	chunkCount, err := indexDocuments(ctx, store, cfg, docs)
	if err != nil {
		return 0, 0, err
	}

	docCount := len(docs)

	// 5. 保存统计信息
	if err := writeCollectionStats(ctx, cfg, docCount, chunkCount); err != nil {
		return 0, 0, fmt.Errorf("更新 Chroma collection 元数据失败: %w", err)
	}

	return docCount, chunkCount, nil
}
