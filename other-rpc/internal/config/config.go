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
	Chroma       ChromaConf    `json:",optional"`
	Memory       MemoryConf    `json:",optional"`
	LLM          LLMConf       `json:",optional"`
	Embedding    EmbeddingConf `json:",optional"`
}

type ChromaConf struct {
	URL        string `json:",default=http://127.0.0.1:8000"`
	Collection string `json:",default=chenaqi_knowledge"`
}

type MemoryConf struct {
	WindowTurns int `json:",default=5"`
	SessionTTL  int `json:",default=1800"`
}

type LLMConf struct {
	Provider string `json:",default=openai"`
	Model    string `json:",optional"`
	BaseURL  string `json:",optional"`
	APIKey   string `json:",optional" env:"LLMAPIKEY"`
}

type EmbeddingConf struct {
	Provider  string `json:",optional"`
	Model     string `json:",optional"`
	BaseURL   string `json:",optional"`
	APIKey    string `json:",optional"`
	Dimension int    `json:",optional"`
}
