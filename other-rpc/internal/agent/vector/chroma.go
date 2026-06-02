package vector

import (
	"context"
	"fmt"
	"strings"

	"other-rpc/internal/config"

	chromav2 "github.com/amikos-tech/chroma-go/pkg/api/v2"

	"github.com/tmc/langchaingo/embeddings"
)

func chromaBaseURL(cfg config.KnowledgeBaseConf) string {
	if cfg.Chroma.URL != "" {
		return strings.TrimRight(cfg.Chroma.URL, "/")
	}
	return "http://127.0.0.1:8000"
}

func chromaCollection(cfg config.KnowledgeBaseConf) string {
	if cfg.Chroma.Collection != "" {
		return cfg.Chroma.Collection
	}
	return "chenaqi_knowledge"
}

func newChromaClient(cfg config.KnowledgeBaseConf) (chromav2.Client, error) {
	return chromav2.NewHTTPClient(chromav2.WithBaseURL(chromaBaseURL(cfg)))
}

func ensureChromaReady(ctx context.Context, cfg config.KnowledgeBaseConf) error {
	client, err := newChromaClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Heartbeat(ctx); err != nil {
		return fmt.Errorf("Chroma 服务不可用 (%s): %w", chromaBaseURL(cfg), err)
	}
	return nil
}

func resetChromaCollection(ctx context.Context, cfg config.KnowledgeBaseConf) error {
	if err := ensureChromaReady(ctx, cfg); err != nil {
		return err
	}

	client, err := newChromaClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.DeleteCollection(ctx, chromaCollection(cfg)); err != nil {
		// collection 不存在时忽略，后续会重新创建
	}
	return nil
}

func openChromaStore(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (*ChromaStore, error) {
	client, err := newChromaClient(cfg)
	if err != nil {
		return nil, err
	}

	collection, err := client.GetOrCreateCollection(ctx, chromaCollection(cfg))
	if err != nil {
		client.Close()
		return nil, err
	}

	return &ChromaStore{
		collection: collection,
		embedder:   embedder,
	}, nil
}

func readCollectionStats(ctx context.Context, cfg config.KnowledgeBaseConf) (docCount, chunkCount int, err error) {
	if err := ensureChromaReady(ctx, cfg); err != nil {
		return 0, 0, err
	}

	client, err := newChromaClient(cfg)
	if err != nil {
		return 0, 0, err
	}
	defer client.Close()

	collection, err := client.GetCollection(ctx, chromaCollection(cfg))
	if err != nil {
		return 0, 0, fmt.Errorf("Chroma collection 不存在，请先 Build: %w", err)
	}

	count, err := collection.Count(ctx)
	if err != nil {
		return 0, 0, err
	}
	chunkCount = count

	if md := collection.Metadata(); md != nil {
		if v, ok := md.GetInt("doc_count"); ok {
			docCount = int(v)
		}
	}
	return docCount, chunkCount, nil
}

func writeCollectionStats(ctx context.Context, cfg config.KnowledgeBaseConf, docCount, chunkCount int) error {
	client, err := newChromaClient(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	collection, err := client.GetCollection(ctx, chromaCollection(cfg))
	if err != nil {
		return err
	}

	return collection.ModifyMetadata(ctx, chromav2.NewMetadataFromMap(map[string]any{
		"doc_count":   int64(docCount),
		"chunk_count": int64(chunkCount),
	}))
}
