# Agent 模块分层

当前 `other-rpc` 的 Agent 按以下层次组织：

- `llm/`：只负责 ChatModel 初始化。
- `embedding/`：只负责 Embedding 初始化。
- `vector/`：负责知识库文档加载、切分、向量化，以及 **Chroma** 向量库读写。
- `rag/`：负责检索问答链封装（`Ask`）。
- `agent.go`：最终 Agent 运行时，聚合以上模块并对外提供 `Chat` / `Build`。

## 调用链路

1. RPC `Chat` 请求进入 `internal/logic/chat_logic.go`。
2. `ChatLogic` 调用 `svcCtx.Agent.Chat(...)`。
3. `Agent` 调用 `rag.QA.Ask(...)`。
4. `rag` 内部执行 RetrievalQA（`LLM + Retriever`）。
5. `Retriever` 来自 Chroma 中已构建的 collection。

## 向量索引生命周期

1. **启动 Chroma 服务**（默认 `http://127.0.0.1:8000`）。
2. **构建（一次性 / 文档更新后）**：调用 gRPC `Agent.Build`，从 `DataPath` 读取文档，切分并向量化，写入 Chroma collection。
3. **启动 agent.rpc**：仅连接 Chroma 加载已有 collection，**不会**重新扫描文档或向量化。
4. **问答**：调用 gRPC `Agent.Chat`。

## 配置

```yaml
KnowledgeBase:
  DataPath: ./data/knowledge
  Chroma:
    URL: http://127.0.0.1:8000
    Collection: chenaqi_knowledge
```

## 初始化入口

`internal/svc/service_context.go`：

- `agent.New(c.KnowledgeBase)` 初始化完整 Agent；
- 启动时连接 Chroma 并验证 collection 中已有向量数据。
