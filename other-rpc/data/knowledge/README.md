# 个人知识库

将 `.txt`、`.md`、`.markdown` 文件放在此目录，启动 `agent.rpc` 时会自动切分、向量化并建立索引。

问答请调用 gRPC `Agent.Chat`，请求体仅包含 `question` 字段。
