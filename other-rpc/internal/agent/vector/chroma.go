package vector

import (
	"context"
	"fmt"
	"other-rpc/internal/config"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"

	"github.com/tmc/langchaingo/embeddings"
)

func NewChromaClient(ctx context.Context, cfg config.KnowledgeBaseConf) (chroma.Client, error) {
	client, err := chroma.NewHTTPClient(chroma.WithBaseURL(cfg.Chroma.URL))
	if err != nil {
		return nil, err
	}

	if err = client.Heartbeat(ctx); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func resetChromaCollection(ctx context.Context, cfg config.KnowledgeBaseConf) error {
	client, err := NewChromaClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.DeleteCollection(ctx, cfg.Chroma.Collection); err != nil {
		// collection 不存在时忽略，后续会重新创建
	}
	return nil
}

func openChromaStore(ctx context.Context, cfg config.KnowledgeBaseConf, embedder embeddings.Embedder) (*ChromaStore, error) {
	client, err := NewChromaClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	collection, err := client.GetOrCreateCollection(ctx, cfg.Chroma.Collection)
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
	client, err := NewChromaClient(ctx, cfg)
	if err != nil {
		return 0, 0, err
	}
	defer client.Close()

	collection, err := client.GetCollection(ctx, cfg.Chroma.Collection)
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
	client, err := NewChromaClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	collection, err := client.GetCollection(ctx, cfg.Chroma.Collection)
	if err != nil {
		return err
	}

	return collection.ModifyMetadata(ctx, chroma.NewMetadataFromMap(map[string]any{
		"doc_count":   int64(docCount),
		"chunk_count": int64(chunkCount),
	}))
}
