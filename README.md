# Chenaqi Backend · 部署指南

基于 [go-zero](https://github.com/zeromicro/go-zero) 的博客系统后端，采用 **HTTP 网关 + gRPC 微服务** 架构。本目录提供一键 Docker Compose 部署，涵盖基础设施、核心业务与 AI 知识库服务。

注意:该README.md内容当前由ai生成，仅供参考

---

## 架构概览

```
                    ┌─────────────────────────────────────────┐
                    │              chenaqi-net                 │
                    │                                         │
  Browser / Client  │   ┌─────────┐      gRPC      ┌────────┐ │
       ────────────►│   │ gateway │ ─────────────► │core-rpc│ │
       HTTP :9000   │   │  :9000  │   (etcd 发现)  │ :8080  │ │
                    │   └────┬────┘                └───┬────┘ │
                    │        │                         │      │
                    │        │ Redis                   │ MySQL │
                    │        ▼                         ▼      │
                    │   ┌─────────┐                ┌────────┐  │
                    │   │  redis  │                │ mysql  │  │
                    │   └─────────┘                └────────┘  │
                    │                                         │
                    │   ┌──────────┐   Chroma   ┌──────────┐  │
                    │   │ other-rpc│ ◄────────► │  chroma  │  │
                    │   │  :8081   │            │  :8000   │  │
                    │   └────┬─────┘            └──────────┘  │
                    │        │ Ollama Embedding                │
                    │        ▼                                 │
                    │   ┌─────────┐      ┌─────────┐          │
                    │   │  ollama │      │  etcd   │          │
                    │   │ :11434  │      │  :2379  │          │
                    │   └─────────┘      └─────────┘          │
                    └─────────────────────────────────────────┘
```

| 服务 | 端口 | 职责 |
|------|------|------|
| **gateway** | 9000 | HTTP API 入口：路由聚合、JWT 鉴权、验证码、文件上传、静态资源 |
| **core-rpc** | 8080 | 核心业务 gRPC：用户、文章、分类、评论、点赞互动 |
| **other-rpc** | 8081 | Agent gRPC：基于 RAG 的个人知识库问答 |
| **mysql** | 3306 | 业务数据持久化 |
| **redis** | 6379 | 验证码缓存、Token 黑名单 |
| **etcd** | 2379 | 服务注册与发现 |
| **chroma** | 8000 | 向量数据库 |
| **ollama** | 11434 | 本地 Embedding（`nomic-embed-text`） |

---

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.25 |
| 框架 | go-zero（REST + zrpc） |
| RPC | gRPC + Protocol Buffers |
| ORM | GORM + MySQL 8.4 |
| 缓存 | Redis 7 |
| 服务发现 | etcd 3.6 |
| AI / RAG | langchaingo、Chroma、Ollama、智谱 GLM |
| 代码生成 | goctl |
| 容器 | Docker Compose |

---

## 快速开始

### 环境要求

- Docker Desktop 或 Docker Engine + Docker Compose v2
- 首次启动需联网（拉取镜像、`ollama pull nomic-embed-text`）
- **other-rpc** 对话能力需配置智谱 API Key（见 [Agent 服务](#agent-服务-other-rpc)）

### 一键部署

```bash
cd backend/deploy
docker compose up -d --build
```

查看运行状态：

```bash
docker compose ps
docker compose logs -f gateway
```

停止并清理：

```bash
docker compose down        # 保留数据卷
docker compose down -v     # 同时删除数据卷
```

### 验证

| 检查项 | 命令 / 地址 |
|--------|-------------|
| Gateway | `curl http://localhost:9000` |
| MySQL 初始化 | 首次启动自动执行 `core-rpc/desc/sql/*.sql` |
| Ollama 模型 | `ollama-pull` 容器首次自动拉取 `nomic-embed-text` |
| 前端联调 | 前端加入 `chenaqi-net` 网络，代理到 `http://gateway:9000` |

---

## 项目结构

```
backend/
├── gateway/          # HTTP API 网关（go-zero REST）
│   ├── desc/         # API 定义（goctl 生成入口 index.api）
│   ├── etc/          # gateway.yaml / gateway-docker.yaml
│   └── static/       # 头像、博客图片
├── core-rpc/         # 核心业务 gRPC
│   ├── desc/
│   │   ├── core.proto
│   │   └── sql/      # MySQL 初始化脚本
│   └── etc/          # core.yaml / core-docker.yaml
├── other-rpc/        # Agent 知识库 gRPC
│   ├── desc/agent.proto
│   ├── data/knowledge/   # 知识库文档（Docker 只读挂载）
│   └── internal/agent/     # RAG 实现（llm / embedding / vector / rag）
└── deploy/           # ← 当前目录：Docker Compose 与部署文档
    ├── docker-compose.yml
    ├── mysql.md / redis.md / Etcd.md
    ├── mq/kafka/           # Kafka 可选独立部署
    └── es_kibana/          # Elasticsearch + Kibana 可选独立部署
```

---

## 业务模块

### Gateway HTTP API

通过 `gateway/desc/index.api` 聚合，主要模块如下：

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| login | `/login/*` | 注册、登录、邮箱验证码、重置密码、登出 |
| user | `/user/*` | 用户资料、我的文章、点赞文章（需鉴权） |
| article | `/article/*` | 文章 CRUD、列表、搜索、按分类查询 |
| category | `/category/*` | 分类创建 / 删除 / 列表 |
| comment | `/comment/*` | 评论、回复、删除、列表 |
| interaction | `/interaction/*` | 浏览量、点赞 / 取消、点赞状态 |
| upload | `/upload/*` | 头像、博客图片上传 |

> Agent 的 HTTP 接口已在 `desc/agent/agent.api` 中定义，尚未接入 `index.api`，当前需直接调用 **other-rpc gRPC**。

### core-rpc 核心能力

- **用户**：注册、邮箱 / 用户名登录、资料读写、登出与 Token 黑名单
- **文章**：创建、编辑、软删除、详情、分页、分类筛选、关键词搜索
- **分类**：增删查（软删除）
- **评论**：一级评论、二级回复、软删除、列表
- **互动**：文章 / 评论点赞、浏览量统计、点赞状态查询

### Agent 服务（other-rpc）

基于 **langchaingo RAG** 的个人知识库问答，gRPC 接口：

| 方法 | 说明 |
|------|------|
| `Build` | 读取 `data/knowledge/` 文档 → 切分 → Ollama 向量化 → 写入 Chroma |
| `Chat` | 检索增强问答（LLM 默认智谱 `glm-4v-flash`） |
| `Test` | 连通性测试 |

**使用前注意：**

1. 在 `other-rpc/etc/agent.yaml`（或 Docker 版 `agent-docker.yaml`）中配置有效的 **LLM API Key**
2. 首次使用前需调用一次 `Build` 建立向量索引
3. 知识库文档更新后需重新 `Build`（启动时不会自动重建已有 collection）
4. 默认知识库目录：`other-rpc/data/knowledge/`

---

## 数据库

Docker 首次启动时，MySQL 自动执行 `core-rpc/desc/sql/` 下的初始化脚本：

| 表 | 说明 |
|----|------|
| `user` | 用户（角色 admin / user / guest，软删除） |
| `category` | 文章分类 |
| `article` | 文章及浏览 / 点赞 / 评论计数 |
| `comment` | 一级 / 二级评论 |
| `interaction_like` | 文章、评论点赞记录 |
| `token_blacklist` | Refresh Token 黑名单 |

运行时 core-rpc 亦通过 GORM `AutoMigrate` 同步 entity 表结构。

---

## 本地开发（GoLand）

不跑 Docker 业务容器、仅在 IDE 中调试时，需先启动基础设施：

```bash
# 参考本目录文档单独启动，或使用 compose 只启中间件：
docker compose up -d etcd mysql redis chroma ollama ollama-pull
```

本地配置文件：

| 服务 | 配置文件 |
|------|----------|
| gateway | `gateway/etc/gateway.yaml` |
| core-rpc | `core-rpc/etc/core.yaml` |
| other-rpc | `other-rpc/etc/agent.yaml` |

启动顺序：

```bash
# 1. core-rpc
cd backend/core-rpc && go run core.go -f etc/core.yaml

# 2. other-rpc（可选）
cd backend/other-rpc && go run agent.go -f etc/agent.yaml

# 3. gateway
cd backend/gateway && go run gateway.go -f etc/gateway.yaml
```

> **注意**：本地 `core.yaml` 中 MySQL 端口默认为 `3307`，Docker Compose 映射为 `3306`，请按实际环境修改。

---

## 代码生成

各子模块目录下执行：

```bash
# gateway：由 API 定义生成 handler / logic / routes
cd gateway && make gateway-code

# core-rpc：由 proto 生成 gRPC 代码
cd core-rpc && make core-protos

# other-rpc：由 proto 生成 gRPC 代码
cd other-rpc && make agent-protos
```

---

## 网络与前端联调

Compose 会创建名为 **`chenaqi-net`** 的 Docker 网络。同一网络内的容器可通过服务名互相访问（如 `gateway:9000`）。

前端容器联调示例（在 `frontend/deploy` 目录）：

```bash
# 1. 先启后端（创建 chenaqi-net）
cd backend/deploy && docker compose up -d --build

# 2. 再启前端（加入 chenaqi-net）
cd frontend/deploy
docker compose -f docker-compose.network.yml up -d --build
```

前端仅容器、Gateway 在宿主机（GoLand）运行时：

```bash
docker compose -f docker-compose.local.yml up -d --build
```

---

## 可选组件

本仓库主 Compose **不包含** 以下组件，可按需独立部署：

| 目录 | 说明 |
|------|------|
| [mysql.md](./mysql.md) | 单独启动 MySQL |
| [redis.md](./redis.md) | 单独启动 Redis |
| [Etcd.md](./Etcd.md) | 单独启动 etcd |
| [mq/kafka/](./mq/kafka/) | Kafka + Kafka UI |
| [es_kibana/](./es_kibana/) | Elasticsearch + Kibana |

---

## 常见问题

<details>
<summary><strong>首次启动较慢？</strong></summary>

`ollama-pull` 需下载 `nomic-embed-text` 模型，视网络情况可能需要数分钟。可通过 `docker compose logs ollama-pull` 查看进度。
</details>

<details>
<summary><strong>other-rpc 启动后 Chat 无结果？</strong></summary>

需先调用 gRPC `Build` 建立知识库索引，并确认 `agent.yaml` 中 LLM API Key 有效、Chroma / Ollama 均已就绪。
</details>

<details>
<summary><strong>前端容器报 network chenaqi-net not found？</strong></summary>

使用 `docker-compose.network.yml` 前须先启动本目录的后端 Compose，以创建 `chenaqi-net` 网络。仅联调宿主机 Gateway 时使用 `docker-compose.local.yml`。
</details>

<details>
<summary><strong>如何更新知识库文档？</strong></summary>

修改 `other-rpc/data/knowledge/` 下的文件后，重新调用 `Agent.Build`。Chroma 中已有 collection 时，other-rpc 重启不会自动重建索引。
</details>

---

## License

Private project — 部署配置仅供本项目使用。
