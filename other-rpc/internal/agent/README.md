# Agent 模块分层

当前 `other-rpc` 的 Agent 按以下层次组织：

- `llm/`：只负责 ChatModel 初始化。
- `embedding/`：只负责 Embedding 初始化。
- `kb/`：负责知识库文档加载、切分、向量化、索引构建。
- `rag/`：负责检索问答链封装（`Ask`）。
- `runtime/`：最终 Agent 运行时，聚合以上模块并对外提供 `Chat`。

## 调用链路

1. RPC `Chat` 请求进入 `internal/logic/chat_logic.go`。
2. `ChatLogic` 调用 `svcCtx.Agent.Chat(...)`。
3. `runtime.Agent` 调用 `rag.QA.Ask(...)`。
4. `rag` 内部执行 RetrievalQA（`LLM + Retriever`）。
5. `Retriever` 来自 `kb.Build(...)` 构建的向量索引。

## 初始化入口

`internal/svc/service_context.go`：

- `runtime.New(c.KnowledgeBase)` 初始化完整 Agent；
- 启动时自动加载 `KnowledgeBase.DataPath` 下文档并建索引。
