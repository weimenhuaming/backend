# Gateway 统一响应格式改造说明

本文档说明 gateway 服务如何从「在 `.api` 和 logic 中手动拼装 `code/msg/data`」迁移到 **go-zero 模版 + 统一 response 包** 的方案。

参考：[go-zero 模板自定义](https://go-zero.dev/zh-cn/reference/customization/template/)

---

## 目标响应格式

所有 API 统一返回：

```json
{
  "code": 0,
  "msg": "OK",
  "data": {}
}
```

- `code: 0` 表示成功
- `data` 为实际业务数据
- 业务错误通过 `code` + `msg` 表达，HTTP 状态码仍为 200

---

## 整体思路

不改 Makefile 的生成方式（`-style go_zero`），只把「统一响应」接入 go-zero 的模版生成链路：

```
.api 只定义业务数据
  → logic 返回 data + error
  → handler 模版统一 response.Response(w, data, err)
  → 前端收到 { code, msg, data }
```

---

## 改造步骤

### 1. 新增统一响应包

路径：`internal/response/response.go`

Handler 统一调用：

```go
response.Response(w, data, err)
```

Logic 中的业务错误使用：

```go
return nil, response.NewError(400, "参数错误")
```

核心逻辑：

| 场景 | code | msg |
|------|------|-----|
| 成功 | `0` | `"OK"` |
| 业务错误 | 自定义（400/401/403/500 等） | 错误信息 |
| 未知 error | `-1` | `err.Error()` |

### 2. 自定义 goctl Handler 模版

路径：`.goctl/api/handler.tpl`

将默认生成的：

```go
if err != nil {
    httpx.ErrorCtx(r.Context(), w, err)
} else {
    httpx.OkJsonCtx(r.Context(), w, resp)
}
```

替换为：

```go
response.Response(w, resp, err)
// 或无返回数据的接口：
response.Response(w, nil, err)
```

模版仅在存在请求体解析时才导入 `httpx`，避免无请求接口产生 unused import。

### 3. Makefile 指向自定义模版

`gateway/Makefile`：

```makefile
gateway-code:
	goctl api go -api desc/index.api -dir . -style go_zero -home .goctl
```

生成命令：

```bash
cd gateway
make gateway-code
```

**注意**：必须使用 `-style go_zero`，文件名才会是 `login_handler.go` 这种下划线风格。

### 4. 精简 `.api` 定义

响应类型中不再手写 `Code`、`Msg`，只保留业务数据结构。

**改造前：**

```api
type LoginResp {
    Code int       `json:"code"`
    Msg  string    `json:"msg"`
    Data LoginData `json:"data"`
}

post /login (LoginReq) returns (LoginResp)
```

**改造后：**

```api
type LoginData {
    Id           uint64 `json:"id"`
    Name         string `json:"name"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    // ...
}

post /login (LoginReq) returns (LoginData)
```

无返回数据的接口不写 `returns`：

```api
post /register (RegisterReq)
```

### 5. 调整 Logic 层

Logic 不再返回带 `code/msg` 的 Resp 结构体。

| 场景 | 改造前 | 改造后 |
|------|--------|--------|
| 成功有数据 | `return &types.LoginResp{Code: 200, Data: ...}, nil` | `return &types.LoginData{...}, nil` |
| 业务错误 | `return &types.LoginResp{Code: 400, Msg: "..."}, nil` | `return nil, response.NewError(400, "...")` |
| 成功无数据 | `return &types.RegisterResp{Code: 200}, nil` | `return nil` |
| 函数签名（无返回） | `(resp *types.RegisterResp, err error)` | `error` |

示例（删除分类）：

```go
func (l *DeleteCategoryLogic) DeleteCategory(req *types.DeleteCategoryReq) error {
    if role != "admin" {
        return response.NewError(403, "非管理员，没有权限执行")
    }
    _, err := l.svcCtx.Core.DeleteCategory(l.ctx, &core_client.DeleteCategoryReq{Id: req.Id})
    if err != nil {
        return response.NewError(500, err.Error())
    }
    return nil
}
```

### 6. 对齐 Makefile 的下划线命名

项目使用 `-style go_zero`，生成文件命名规则：

| 类型 | 命名示例 |
|------|----------|
| Handler | `login_handler.go` |
| Logic | `delete_category_logic.go` |
| Middleware | `auth_middleware.go` |
| ServiceContext | `service_context.go` |

不要使用 `-style gozero`（无下划线），否则会产生 `loginhandler.go` 等同功能重复文件。

真实业务实现应保留在下划线文件中，删除 goctl 生成的空壳副本。

### 7. 特殊 Handler 手动保留

登录相关接口有 Cookie / Header 逻辑，不能只用模版默认生成，需在下划线 handler 中手动处理：

| 文件 | 特殊逻辑 |
|------|----------|
| `login_handler.go` | SetCookie(refresh_token)、Authorization 头、响应前清空 RefreshToken |
| `email_login_handler.go` | 同上 |
| `logout_handler.go` | 从 Cookie 读取 refresh_token 再调用 logic |

示例（登录 handler 核心片段）：

```go
resp, err := l.Login(&req)
if err != nil {
    response.Response(w, nil, err)
    return
}

http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: resp.RefreshToken, ...})
resp.RefreshToken = ""
w.Header().Set("Authorization", resp.AccessToken)
response.Response(w, resp, nil)
```

其余普通 handler 由模版自动生成即可。

---

## 数据流对比

### 改造前

```
.api 定义 Code/Msg/Data
  → logic 手动组装完整响应
  → handler 直接 httpx.OkJsonCtx(w, resp)
  → 前端: { "code": 200, "msg": "登录成功", "data": {...} }
```

### 改造后

```
.api 只定义业务数据
  → logic 返回 data + error
  → handler 模版统一 response.Response(w, data, err)
  → 前端: { "code": 0, "msg": "OK", "data": {...} }
```

---

## 涉及文件一览

| 路径 | 作用 |
|------|------|
| `internal/response/response.go` | 统一响应封装 |
| `.goctl/api/handler.tpl` | 自定义 handler 生成模版 |
| `Makefile` | 代码生成入口 |
| `desc/**/*.api` | API 定义（已去掉 Code/Msg） |
| `internal/handler/**/*_handler.go` | HTTP 入口 |
| `internal/logic/**/*_logic.go` | 业务逻辑 |
| `internal/middleware/auth_middleware.go` | 认证中间件（保留真实实现） |
| `internal/svc/service_context.go` | 服务上下文（保留真实实现） |

---

## 新增接口流程

1. 在 `desc/` 下对应 `.api` 中定义请求和业务数据类型
2. 执行 `make gateway-code` 生成 handler / logic / types
3. 在 `internal/logic/` 对应文件中实现业务逻辑
4. 错误使用 `response.NewError(code, msg)` 返回
5. 一般无需修改 handler；若有 Cookie/Header 等特殊需求，再手动改 handler

---

## 对前端的影响

| 项目 | 改造前 | 改造后                 | 是否影响        |
|------|--------|---------------------|-------------|
| 响应结构 | `{code, msg, data}` | `{code, msg, data}` | 无           |
| `data` 字段内容 | 业务数据 | 业务数据                | 无           |
| 成功 `code` | `200` | `200`               | 无           |
| 成功 `msg` | `"登录成功"` 等 | `"OK"`              | 若展示 msg 需调整 |
| 业务错误 `code` | `400/401/500` | `400/401/500`       | 无           |
| HTTP 状态码 | 业务错误也是 200 | 仍是 200              | 无           |

前端成功判断建议改为：

```javascript
if (res.code === 0) {
  // 成功，使用 res.data
}
```

过渡期兼容写法：

```javascript
if (res.code === 0 || res.code === 200) {
  // 兼容新旧后端
}
```

---

## 常见问题

### Q: 为什么不在 `.api` 里继续写 Code/Msg？

go-zero 官方推荐由 handler 层统一包装响应。`.api` 只描述业务数据结构，避免每个接口重复定义相同字段，logic 也更纯粹。

### Q: 重新生成会覆盖我的 logic 吗？

`make gateway-code` 对已存在的 logic 文件会 **跳过生成**（`exists, ignored generation`）。已实现的 logic 不会被覆盖；新增的接口才会生成 stub。

### Q: 重新生成会覆盖 handler 吗？

同样会跳过已存在的 handler。若需让模版重新生成某个 handler，需先删除该文件再执行 `make gateway-code`。

### Q: 参数校验失败（httpx.Parse）返回什么？

请求体解析失败仍走 go-zero 默认的 `httpx.ErrorCtx`，不会进入 `response.Response`。如需统一格式，可进一步改造模版或中间件。

---

## 参考

- [go-zero 模板自定义](https://go-zero.dev/zh-cn/reference/customization/template/)
- [修改模板前后响应体对比](https://go-zero.dev/zh-cn/reference/customization/template/#%E4%BF%AE%E6%94%B9%E6%A8%A1%E6%9D%BF%E5%89%8D%E5%90%8E%E5%93%8D%E5%BA%94%E4%BD%93%E5%AF%B9%E6%AF%94)
