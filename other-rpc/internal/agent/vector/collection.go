package vector

import (
	"context"
	"fmt"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	chromaembed "github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// Collection 基于 Chroma v2 API 的向量存储，实现 langchaingo VectorStore。
type Collection struct {
	Name       string
	DocCount   int
	ChunkCount int
	Count      int
	collection chroma.Collection
	embedder   embeddings.Embedder
}

// NewCollection 新建知识库
func (c *Chroma) NewCollection(ctx context.Context, name string, embedder embeddings.Embedder) (*Collection, error) {
	collection, err := c.client.GetOrCreateCollection(ctx, name, collectionCreateOptions()...)
	if err != nil {
		return nil, err
	}

	return &Collection{
		Name:       name,
		collection: collection,
		embedder:   embedder,
	}, nil
}

func (s *Collection) AddDocuments(ctx context.Context, docs []schema.Document, _ ...vectorstores.Option) ([]string, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	texts := make([]string, len(docs))
	metas := make([]chroma.DocumentMetadata, len(docs))
	for i, doc := range docs {
		texts[i] = doc.PageContent
		metas[i] = toDocumentMetadata(doc.Metadata)
	}

	vectors, err := s.embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}
	chromaVectors := make([]chromaembed.Embedding, len(vectors))
	for i, vec := range vectors {
		chromaVectors[i] = chromaembed.NewEmbeddingFromFloat32(vec)
	}

	if err := s.collection.Add(ctx,
		chroma.WithIDGenerator(chroma.NewUUIDGenerator()),
		chroma.WithTexts(texts...),
		chroma.WithEmbeddings(chromaVectors...),
		chroma.WithMetadatas(metas...),
	); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Collection) SimilaritySearch(ctx context.Context, query string, numDocuments int, _ ...vectorstores.Option) ([]schema.Document, error) {
	if numDocuments <= 0 {
		numDocuments = 4
	}

	queryVector, err := s.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.Query(ctx,
		chroma.WithQueryEmbeddings(chromaembed.NewEmbeddingFromFloat32(queryVector)),
		chroma.WithNResults(numDocuments),
		chroma.WithIncludeQuery(chroma.IncludeDocuments, chroma.IncludeMetadatas, chroma.IncludeDistances),
	)
	if err != nil {
		return nil, err
	}

	docGroups := result.GetDocumentsGroups()
	metaGroups := result.GetMetadatasGroups()
	distGroups := result.GetDistancesGroups()
	if len(docGroups) == 0 {
		return nil, nil
	}

	out := make([]schema.Document, 0, len(docGroups[0]))
	for i, doc := range docGroups[0] {
		item := schema.Document{PageContent: doc.ContentString()}
		if len(metaGroups) > 0 && i < len(metaGroups[0]) && metaGroups[0][i] != nil {
			item.Metadata = fromDocumentMetadata(metaGroups[0][i])
		}
		if len(distGroups) > 0 && i < len(distGroups[0]) {
			item.Score = float32(1.0 - float64(distGroups[0][i]))
		}
		out = append(out, item)
	}
	return out, nil
}

func toDocumentMetadata(meta map[string]any) chroma.DocumentMetadata {
	if len(meta) == 0 {
		return chroma.NewDocumentMetadata()
	}
	attrs := make([]*chroma.MetaAttribute, 0, len(meta))
	for key, value := range meta {
		switch v := value.(type) {
		case string:
			attrs = append(attrs, chroma.NewStringAttribute(key, v))
		case int:
			attrs = append(attrs, chroma.NewIntAttribute(key, int64(v)))
		case int32:
			attrs = append(attrs, chroma.NewIntAttribute(key, int64(v)))
		case int64:
			attrs = append(attrs, chroma.NewIntAttribute(key, v))
		case float64:
			attrs = append(attrs, chroma.NewFloatAttribute(key, v))
		case bool:
			attrs = append(attrs, chroma.NewBoolAttribute(key, v))
		default:
			attrs = append(attrs, chroma.NewStringAttribute(key, fmt.Sprint(v)))
		}
	}
	return chroma.NewDocumentMetadata(attrs...)
}

func fromDocumentMetadata(meta chroma.DocumentMetadata) map[string]any {
	if meta == nil {
		return nil
	}
	out := make(map[string]any)
	if source, ok := meta.GetString("source"); ok {
		out["source"] = source
	}
	return out
}
