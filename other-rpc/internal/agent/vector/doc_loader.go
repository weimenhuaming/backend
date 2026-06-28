package vector

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"other-rpc/internal/config"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

var supportedExts = map[string]struct{}{
	".txt":      {},
	".md":       {},
	".markdown": {},
}

func LoadDocumentsFromDir(dir string) ([]schema.Document, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, nil
	}

	var docs []schema.Document
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if _, ok := supportedExts[ext]; !ok {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.TrimSpace(string(content))
		if text == "" {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		docs = append(docs, schema.Document{
			PageContent: text,
			Metadata: map[string]any{
				"source": rel,
			},
		})
		return nil
	})
	return docs, err
}

// 切割数据库文档，存入向量数据库
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
