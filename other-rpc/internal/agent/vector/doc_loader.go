package vector

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/schema"
)

var supportedExts = map[string]struct{}{
	".txt": {}, ".md": {}, ".markdown": {},
}

// LoadDocumentsFromDir 从目录加载文本/Markdown 为 Document。
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
