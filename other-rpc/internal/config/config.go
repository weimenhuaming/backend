package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	KnowledgeBase KnowledgeBaseConf
}

type KnowledgeBaseConf struct {
	DataPath     string        `json:",default=./data/knowledge"`
	TopK         int           `json:",default=4"`
	ChunkSize    int           `json:",default=800"`
	ChunkOverlap int           `json:",default=100"`
	LLM          LLMConf       `json:",optional"`
	Embedding    EmbeddingConf `json:",optional"`
}

type LLMConf struct {
	Provider string `json:",default=openai"`
	Model    string `json:",optional"`
	BaseURL  string `json:",optional"`
	APIKey   string `json:",optional"`
}

type EmbeddingConf struct {
	Provider string `json:",default=ollama"`
	Model    string `json:",optional"`
	BaseURL  string `json:",optional"`
	APIKey   string `json:",optional"`
}
