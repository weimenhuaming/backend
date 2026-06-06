# 后端部署

在本目录下启动全部后端服务（基础设施 + RPC + Gateway）：

```bash
docker compose up -d --build
```

前端请单独在 `frontend/deploy` 目录部署。

## 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| gateway | 9000 | API 网关 |
| core-rpc | 8080 | 核心业务 RPC |
| other-rpc | 8081 | Agent / 知识库 RPC |
| chroma | 8000 | 向量数据库 |
| ollama | 11434 | 本地 Embedding（nomic-embed-text） |
| mysql | 3306 | 数据库 |
| redis | 6379 | 缓存 |
| etcd | 2379 | 服务注册发现 |

首次启动时 `ollama-pull` 会自动执行 `ollama pull nomic-embed-text`，无需手动 `docker exec`。

知识库文档更新后，需调用 `Agent.Build` 重新向量化（Chroma 中已有 collection 时 other-rpc 启动不会自动重建索引）。

## 网络

Compose 会创建名为 `chenaqi-net` 的 Docker 网络，供前端容器通过服务名 `gateway` 访问 API。

## 可选组件

- `mysql.md` / `redis.md` / `Etcd.md`：单独启动各基础设施的说明
- `mq/kafka/`：Kafka 独立部署
- `es_kibana/`：Elasticsearch + Kibana 独立部署
