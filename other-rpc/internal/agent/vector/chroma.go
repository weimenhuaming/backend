package vector

import (
	"context"
	"errors"
	"fmt"
	"other-rpc/internal/config"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	chromaembed "github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

var (
	ErrCollectionExists   = errors.New("collection 已存在")
	ErrCollectionNotFound = errors.New("collection 不存在")
)

// Chroma 封装 Chroma HTTP 客户端，复用连接执行 collection 操作。
type Chroma struct {
	cfg    config.KnowledgeBaseConf
	client chroma.Client
}

// NewChroma 创建并校验 Chroma 连接。
func NewChroma(ctx context.Context, cfg config.KnowledgeBaseConf) (*Chroma, error) {
	client, err := chroma.NewHTTPClient(chroma.WithBaseURL(cfg.Chroma.URL))
	if err != nil {
		return nil, err
	}

	if err = client.Heartbeat(ctx); err != nil {
		client.Close()
		return nil, err
	}

	// 结合文章的数量
	topK := cfg.TopK
	if topK <= 0 {
		topK = 4
	}

	return &Chroma{
		cfg:    cfg,
		client: client,
	}, nil
}

// Close 关闭 Chroma 客户端。
func (c *Chroma) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Load 加载默认 collection 的检索器；若 collection 不存在则先 Build。
func (c *Chroma) Load(ctx context.Context, embedder embeddings.Embedder) (vectorstores.Retriever, error) {
	name := c.cfg.Chroma.Collection
	collection, err := c.client.GetCollection(ctx, name)
	if err != nil {
		retriever, _, _, err := c.Build(ctx, name, embedder)
		return retriever, err
	}

	store := &Collection{
		collection: collection,
		embedder:   embedder,
	}
	return vectorstores.ToRetriever(store, c.cfg.TopK), nil
}

// Build 从知识库目录构建指定名称的 collection；若名称已存在则返回 ErrCollectionExists。
func (c *Chroma) Build(ctx context.Context, name string, embedder embeddings.Embedder) (vectorstores.Retriever, int, int, error) {
	if name == "" {
		return vectorstores.Retriever{}, 0, 0, errors.New("collection 名称不能为空")
	}

	// 1. 首先判断是否存在
	if _, err := c.client.GetCollection(ctx, name); err == nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("%w: %s", ErrCollectionExists, name)
	}

	// 2.新建collection
	store, err := c.NewCollection(ctx, name, embedder)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("创建 Chroma collection 失败: %w", err)
	}

	// 3.加载和切割文档，存入向量数据库
	docs, err := LoadDocumentsFromDir(c.cfg.DataPath)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("加载知识库目录失败: %w", err)
	}

	chunkCount, err := indexDocuments(ctx, store, c.cfg, docs)
	if err != nil {
		return vectorstores.Retriever{}, 0, 0, err
	}

	docCount := len(docs)

	// 4.写入元数据
	if err := c.writeCollectionStats(ctx, name, docCount, chunkCount); err != nil {
		return vectorstores.Retriever{}, 0, 0, fmt.Errorf("更新 Chroma collection 元数据失败: %w", err)
	}

	return vectorstores.ToRetriever(store, c.cfg.TopK), docCount, chunkCount, nil
}

// OpenRetriever 打开已有 collection 并返回检索器。
func (c *Chroma) OpenRetriever(ctx context.Context, name string, embedder embeddings.Embedder) (vectorstores.Retriever, error) {
	if name == "" {
		return vectorstores.Retriever{}, errors.New("collection 名称不能为空")
	}

	collection, err := c.client.GetCollection(ctx, name)
	if err != nil {
		return vectorstores.Retriever{}, fmt.Errorf("collection %q 不存在", name)
	}

	store := &Collection{
		collection: collection,
		embedder:   embedder,
	}
	return vectorstores.ToRetriever(store, c.cfg.TopK), nil
}

// DeleteCollection 删除指定名称的 collection。
func (c *Chroma) DeleteCollection(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("collection 名称不能为空")
	}

	if _, err := c.client.GetCollection(ctx, name); err != nil {
		return fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}

	if err := c.client.DeleteCollection(ctx, name); err != nil {
		return fmt.Errorf("删除 collection %q 失败: %w", name, err)
	}
	return nil
}

// ListCollections 返回当前数据库下所有 collection 的概要信息。
func (c *Chroma) ListCollections(ctx context.Context) ([]Collection, error) {
	collections, err := c.client.ListCollections(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Collection, 0, len(collections))
	for _, col := range collections {
		info := Collection{Name: col.Name()}
		if meta := col.Metadata(); meta != nil {
			if docCount, ok := meta.GetInt("doc_count"); ok {
				info.DocCount = int(docCount)
			}
			if chunkCount, ok := meta.GetInt("chunk_count"); ok {
				info.ChunkCount = int(chunkCount)
			}
		}
		count, err := col.Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取 collection %q 数量失败: %w", col.Name(), err)
		}
		info.Count = count
		out = append(out, info)
	}
	return out, nil
}

// collectionCreateOptions 返回创建 Chroma collection 时的客户端选项。
//
// 背景：chroma-go 在 CreateCollection 时若未指定 EmbeddingFunction，会默认初始化
// ONNX 本地嵌入（ort.NewDefaultEmbeddingFunction），需加载 libonnxruntime 和 libstdc++.so.6。
// other-rpc 的 Docker 镜像基于 Alpine，不含 libstdc++.so.6，因此在容器内建库会失败；
//
// 本项目实际不向 Chroma 要嵌入：向量由外部 Embedding 服务在客户端生成，写入/检索时通过
// WithEmbeddings / WithQueryEmbeddings 传入（见 collection.go）。
//
// 因此这里：
//   - WithEmbeddingFunctionCreate(ConsistentHash...)：纯 Go 占位，绕过 ONNX 默认嵌入；
//   - WithDisableEFConfigStorage()：不把占位 EF 配置写入 Chroma 服务端。
func collectionCreateOptions() []chroma.CreateCollectionOption {
	return []chroma.CreateCollectionOption{
		chroma.WithEmbeddingFunctionCreate(chromaembed.NewConsistentHashEmbeddingFunction()),
		chroma.WithDisableEFConfigStorage(),
	}
}

// 写入元数据
func (c *Chroma) writeCollectionStats(ctx context.Context, name string, docCount, chunkCount int) error {
	collection, err := c.client.GetCollection(ctx, name)
	if err != nil {
		return err
	}

	return collection.ModifyMetadata(ctx, chroma.NewMetadataFromMap(map[string]any{
		"doc_count":   int64(docCount),
		"chunk_count": int64(chunkCount),
	}))
}
