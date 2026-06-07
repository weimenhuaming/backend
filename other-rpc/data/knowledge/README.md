# Chenaqi 个人博客 · 项目知识库

本知识库供站内 AI 助手「月忆」检索使用，介绍 **Chenaqi Web** 个人博客全栈项目的定位、架构与功能。

---

## 项目概述

**Chenaqi Web** 是博主「陈阿七」的个人博客与内容管理系统，包含：

- 面向读者的博客阅读、评论、点赞与互动
- 面向作者的后台文章、分类与用户管理
- 基于 RAG 的站内 AI 知识库问答（Agent 服务）

仓库采用前后端分离：`frontend/` 为 Nuxt 4 前端，`backend/` 为 go-zero 微服务后端。

---

## 技术栈

### 前端（`frontend/`）

| 类别 | 选型 |
|------|------|
| 框架 | Nuxt 4 · Vue 3 |
| 状态管理 | Pinia |
| HTTP | `$fetch` / `ofetch` |
| Markdown | marked |
| 样式 | 原生 CSS |
| 运行时 | Node.js 22 |

### 后端（`backend/`）

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.25 |
| 框架 | go-zero（REST + zrpc） |
| RPC | gRPC + Protocol Buffers |
| ORM | GORM + MySQL 8.4 |
| 缓存 | Redis 7 |
| 服务发现 | etcd |
| AI / RAG | langchaingo、Chroma、Ollama Embedding、智谱 GLM |

---

## 系统架构

```
Browser → Nuxt Frontend (:3000)
              │ /api/** 代理
              ▼
         Gateway (:9000)  HTTP 入口、JWT 鉴权
              │ gRPC (etcd 发现)
              ├── core-rpc (:8080)  用户 / 文章 / 分类 / 评论 / 点赞
              └── other-rpc (:8081) Agent 知识库问答
                        │
                        ├── Chroma (:8000)  向量数据库
                        └── Ollama (:11434)  Embedding（nomic-embed-text）
```

| 服务 | 端口 | 职责 |
|------|------|------|
| gateway | 9000 | HTTP API 聚合、鉴权、上传、静态资源 |
| core-rpc | 8080 | 核心业务 gRPC |
| other-rpc | 8081 | Agent gRPC：RAG 知识库问答 |
| mysql | 3306 | 业务数据 |
| redis | 6379 | 验证码、Token 黑名单 |
| chroma | 8000 | 向量存储 |
| ollama | 11434 | 本地向量化 |

---

## 前端功能

### 公开页面

| 路由 | 说明 |
|------|------|
| `/` | 首页 Bento 布局（时钟、日历、音乐、推荐等） |
| `/blog` | 博客列表 |
| `/blog/:id` | 文章详情（Markdown、目录、评论、点赞） |
| `/about` | 关于页 |
| `/agent` | 站内 AI 助手 |
| `/auth/login` | 登录 / 注册 / 重置密码 |

### 用户后台（需登录）

| 路由 | 说明 |
|------|------|
| `/admin/profile` | 个人资料、头像上传 |
| `/admin/likes` | 点赞文章列表 |
| `/admin/my-articles` | 我的文章 |

### 管理员后台（需 admin 角色）

| 路由 | 说明 |
|------|------|
| `/admin/categories` | 分类管理 |
| `/admin/articles` | 全站博客列表 |
| `/admin/article/create` | 新建 / 编辑文章 |

---

## 后端业务模块（Gateway HTTP API）

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| login | `/login/*` | 注册、登录、邮箱验证码、重置密码、登出 |
| user | `/user/*` | 用户资料、我的文章、点赞文章 |
| article | `/article/*` | 文章 CRUD、列表、搜索、分类筛选 |
| category | `/category/*` | 分类增删查 |
| comment | `/comment/*` | 评论、回复、删除、列表 |
| interaction | `/interaction/*` | 浏览量、点赞 / 取消、点赞状态 |
| upload | `/upload/*` | 头像、博客图片上传 |

---

## Agent 知识库服务（other-rpc）

基于 **langchaingo RAG** 的个人知识库问答，主要 gRPC 接口：

| 方法 | 说明 |
|------|------|
| `Build` | 从 `data/knowledge/` 读取文档，切分向量化，创建新 collection（名称不可重复） |
| `SwitchRetriever` | 切换当前问答使用的 collection |
| `ListCollections` | 查看所有 collection 及文档 / 切片统计 |
| `DeleteCollection` | 删除指定 collection |
| `Chat` | 检索增强问答 |
| `Test` | 连通性测试 |

### 知识库目录

默认路径：`other-rpc/data/knowledge/`

将 Markdown、文本等文档放入此目录后，调用 `Build` 并指定 collection 名称即可建立索引。文档更新后需对新名称重新 `Build`，或通过 `SwitchRetriever` 切换到对应 collection。

### 配置示例

```yaml
KnowledgeBase:
  DataPath: ./data/knowledge
  TopK: 4
  ChunkSize: 800
  ChunkOverlap: 100
  Chroma:
    URL: http://127.0.0.1:8000
    Collection: chenaqi_knowledge
  LLM:
    Provider: openai
    Model: glm-4-flash
    BaseURL: https://open.bigmodel.cn/api/paas/v4/
  Embedding:
    Provider: ollama
    Model: nomic-embed-text
    BaseURL: http://127.0.0.1:11434
```

---

## 数据库表（MySQL）

| 表 | 说明 |
|----|------|
| user | 用户（角色 admin / user / guest） |
| category | 文章分类 |
| article | 文章及浏览 / 点赞 / 评论计数 |
| comment | 一级 / 二级评论 |
| interaction_like | 文章、评论点赞 |
| token_blacklist | Refresh Token 黑名单 |

---

## 作者与博客

- **博主昵称**：陈阿七
- **项目名**：Chenaqi Web
- **专注方向**：Go 语言开发、微服务架构、云原生、AI 应用
- **博客定位**：分享技术实践与个人成长，帮助开发者少走弯路
- **AI 助手名称**：小陈（站内 RAG 问答助手）

---

## 部署方式

### Docker 一键部署

```bash
cd backend/deploy
docker compose up -d --build
```

包含 gateway、core-rpc、other-rpc、mysql、redis、etcd、chroma、ollama 等。

### 本地开发

1. 启动中间件：`etcd`、`mysql`、`redis`、`chroma`、`ollama`
2. 依次运行：`core-rpc` → `other-rpc` → `gateway`
3. 前端：`cd frontend && npm run dev`（代理到 Gateway :9000）

---

## 常见问题

**首次使用 Agent 问答无结果？**

需先调用 `Build` 建立向量索引，并配置有效的 LLM API Key，确保 Chroma 与 Ollama 均已就绪。

**如何管理多个知识库？**

使用不同 collection 名称多次 `Build`，通过 `ListCollections` 查看，用 `SwitchRetriever` 切换，`DeleteCollection` 删除不需要的库。

**知识库文档放在哪里？**

`backend/other-rpc/data/knowledge/`，支持 Markdown 等文本格式，由 `doc_loader` 加载后切分入库。
