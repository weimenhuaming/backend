# 个人知识库

将 `.txt`、`.md`、`.markdown` 文件放在此目录。

## 使用流程

1. 启动 Chroma 服务（默认 `http://127.0.0.1:8000`，可用 `deploy/docker-compose.yml` 中的 `chroma` 服务）。
2. 调用 gRPC `Agent.Build` 构建向量索引（数据写入 Chroma collection）。
3. 启动 `agent.rpc`（自动连接 Chroma，**不会**在启动时重新构建）。
4. 调用 gRPC `Agent.Chat` 进行问答。

文档更新后，重新调用 `Build` 即可（会清空并重建对应 collection）。
