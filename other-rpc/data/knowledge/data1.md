# 前言
在使用 ChatGPT、DeepSeek、Claude 等大模型产品时，我们会发现一个共同特点： 当用户发送问题后，AI 的回答并不会等全部生成完成再返回，而是像打字一样逐步输出。

```latex
用户：什么是 Go 的 GC？

AI：
Go 的垃圾回收器（Garbage Collector）采用...

随后
Go 的垃圾回收器（Garbage Collector）采用三色标记法...
```

这种技术称为：**<font style="color:#E746A4;">Streaming Response（流式响应）</font>**，流式响应已经成为大模型应用开发中的标准能力。



:::info
为什么需要流式响应

:::

假设用户发送一个问题：请详细解释 Go 的垃圾回收机制

如果采用传统 HTTP 请求：

```latex
Client
   │
   │ Request
   ▼
Server
   │
   │ 等待 LLM 完整生成
   │
   ▼
Response
```

可能需要：

+ 模型思考：3 秒
+ 生成答案：8 秒

用户需要等待11s，期间页面没有任何反馈，用户体验较差。



:::color4
而流式输出模式：

:::

```latex
Client
   │
   │ Request
   ▼
Server
   │
   ├── Token1
   ├── Token2
   ├── Token3
   ├── Token4
   └── ...
   ▼
Client
```

用户可能：1 秒后就看到内容开始输出，虽然总耗时没有减少，但感知速度大幅提升。

# 一.常见流式传输方案
目前主流方案有两种：

## 1.1 SSE（Server-Sent Events）
:::color5
服务端持续向客户端推送数据，属于是单向

:::

```latex
Browser
    │ HTTP
    ▼
Server
    │
    ├── data: hello
    ├── data: world
    ├── data: ...
```

特点：

+ 基于 HTTP
+ 实现简单
+ 浏览器原生支持
+ 非常适合 AI 问答

## 1.2 WebSocket
:::color3
建立长连接后双向通信。

:::

```latex
Client
   ⇄
WebSocket
   ⇄
Server
```

特点：

+ 双向通信
+ 更灵活
+ 实现复杂度较高

## 1.3 为什么大模型场景更推荐 SSE
AI 聊天场景通常是：

```latex
用户发送问题
      ↓
服务器持续返回答案
```

本质上是：

```latex
Client → Server
Server → Client
```

单向数据流，因此 SSE 更简单，目前很多大模型接口都采用 SSE：

+ OpenAI
+ DeepSeek
+ Claude
+ Gemini

对于绝大多数知识库、RAG、AI 问答系统而言，优先推荐使用 SSE 实现流式输出，因为它足够简单、稳定，并且与主流大模型 API 的设计保持一致。

# 二.SSE 工作原理和实现
浏览器发起请求：

```http
GET /chat/stream
```

服务端保持连接不断开：

```latex
data: 你好

data: 我是 AI 助手

data: 很高兴为你服务
```

浏览器持续接收。

## 2.1 SSE 消息格式
SSE 规范要求：

```latex
data: 内容
```

注意：每条消息后必须有两个换行。

Go 中通常这样写：

```go
fmt.Fprintf(w, "data: %s\n\n", content)
```

## 2.2 Go 实现 SSE
```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", 500)
        return
    }

    messages := []string{
        "你好",
        "我是",
        "一个",
        "流式响应",
        "示例",
    }

    for _, msg := range messages {

        fmt.Fprintf(w, "data: %s\n\n", msg)

        flusher.Flush()

        time.Sleep(time.Second)
    }
}

func main() {
    http.HandleFunc("/stream", streamHandler)
    http.ListenAndServe(":8080", nil)
}
```



浏览器原生支持 EventSource。

```javascript
const eventSource = new EventSource("/stream");

eventSource.onmessage = function(event) {
    console.log(event.data);
};
```

```latex
你好
我是
一个
流式响应
示例
```

:::color3
整体上 Go 服务只是一个流量中转站

:::

## 2.3 Go-Zero 中实现流式响应
```go
func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set(
            "Content-Type",
            "text/event-stream",
        )

        flusher := w.(http.Flusher)

        logic := logic.NewChatLogic(
            r.Context(),
            svcCtx,
        )

        logic.StreamChat(func(content string) {

            fmt.Fprintf(
                w,
                "data: %s\n\n",
                content,
            )

            flusher.Flush()
        })
    }
}
```

```go
func (l *ChatLogic) StreamChat(
    send func(string),
) {

    stream, err := l.llm.Stream(...)

    if err != nil {
        return
    }

    for stream.Next() {

        token := stream.Content()

        send(token)
    }
}
```

## 2.4 LangChainGo 中的流式输出
```go
llm.GenerateContent(
    ctx,
    messages,
    llms.WithStreamingFunc(
        func(
            ctx context.Context,
            chunk []byte,
        ) error {

            fmt.Print(
                string(chunk),
            )

            return nil
        },
    ),
)
```

每生成一个 Token 就会触发一次回调，随后即可通过 SSE 推送给前端。



# 三.具体案例
## 3.1 案例介绍
这里详细介绍一下如何使用SSE做一个流式处理，这里以gin框架为例，原生的http和gozero应该差不多，下面展示的案例，只用看stream处理部分即可

1. 首先第一步，就是处理请求头，告诉它我们是一个流式处理
2. 设置回调函数 func(chunk string) error  （在后续处理中，会返回一个chunk，也就是部分文本，通过这个回调函数，解析成json格式，然后按照SSE的格式通过Flush将缓冲区的数据发给客户端）
3. 最后通知其结束了。

```go
func (ctrl *VolcEngineController) Chat(c *gin.Context) {
    // ...
    
    // 1. 加上请求头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

    // 2. 设置回调函数
	err := Service.VolcEngineSvc.Result(req.Pid, params, func(chunk string) error {
		payload, marshalErr := json.Marshal(map[string]string{"content": chunk})
		if marshalErr != nil {
			return marshalErr
		}
		if _, writeErr := fmt.Fprintf(c.Writer, "data: %s\n\n", payload); writeErr != nil {
			return writeErr
		}
		c.Writer.Flush()
		return nil
	})
	if err != nil {
		errPayload, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", errPayload)
		c.Writer.Flush()
		return
	}

    // 3.最后通知结束了
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
}

```

接着看一下里面具体的逻辑（由于采用的并不是SDK而是curl需新建request，所以我后续还需要在加一次header；还有数据的发送，这里我就直接跳过了，来看它如何读取里面的数据的。）

下面这个代码做的就是从响应中读取SSE格式流数据和解析对应响应

```go
func readOpenAIStream(body io.Reader, callback ChatStreamCallback) error {
	// 1.创建带缓冲的读取器
    reader := bufio.NewReader(body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
            // io.EOF表示正常结束
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read stream failed: %w", err)
		}
        // 过滤无效行
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
    // 2.提取数据，去掉data前缀
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return nil
		}
    // 3.定义和解析json数据
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error,omitempty"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if chunk.Error != nil {
			return fmt.Errorf("api error: %s", chunk.Error.Message)
		}
    // 4.提取内容并回调
		if len(chunk.Choices) == 0 {
			continue
		}
		content := chunk.Choices[0].Delta.Content
		if content == "" {
			continue
		}
		if err := callback(content); err != nil {
			return err
		}
	}
}
```

## 3.2 bufiio.NewReader
`body` 本身已经是 `io.Reader`，但它可能没有缓冲，直接读取效率低且不方便按行读取。`bufio.NewReader` 给它加上了一层缓冲区，提供便捷的 `ReadString('\n')` 方法。

bufiio.NewReader本质上做了两件事：

1. 批量读取：一次系统调用读4096字节到缓冲区
2. 提供一些便利的方法：readString等

:::info
实际上做到了双层缓冲区的一个效果

:::

![](https://cdn.nlark.com/yuque/__mermaid_v3/5349623bdcd692af36bc5c7ff2094af7.svg)

1. 内核 TCP 缓冲区：操作系统管理，自动接收网络数据
2. resp.Body：直接读取内核缓冲区，无应用层缓冲
3. bufio.Reader：应用层缓冲，提供按行读取能力



:::color4
这样做的原因：

:::

1. 首先第一点：无bufio的，比较复杂
2. 可以减少系统调用次数，提高效率

```go
// SSE 返回的数据流
// 网络包 1: "data: {\"content\":\"你\"}\n\n"
// 网络包 2: "data: {\"content\":\"好\"}\n\n"
// 网络包 3: "data: {\"content\":\"世\"}\n\n"
// 网络包 4: "界\"}\n\ndata: [DONE]\n\n"

func compareReadPerformance() {
    // 假设 body 是真实的网络连接
    body := getHttpResponseBody()
    
    // ❌ 无缓冲：4 次系统调用（每行一次）
    start := time.Now()
    buf := make([]byte, 1024)
    for i := 0; i < 4; i++ {
        n, _ := body.Read(buf)  // 系统调用 1,2,3,4
        // 处理...
    }
    fmt.Println("无缓冲:", time.Since(start))
    
    // ✅ 有缓冲：1-2 次系统调用
    start = time.Now()
    body2 := getHttpResponseBody()  // 新的连接
    reader := bufio.NewReader(body2)
    for i := 0; i < 4; i++ {
        line, _ := reader.ReadString('\n')  // 第 1 次触发系统调用（读 4096 字节）
        // 第 2,3,4 次直接从缓冲区返回，无系统调用
    }
    fmt.Println("有缓冲:", time.Since(start))
}

// 输出示例：
// 无缓冲: 2.3ms
// 有缓冲: 0.4ms  (快 5-6 倍)
```

resp.Body 直接连接 TCP socket，没有应用层缓冲区，也不提供按行读取方法。

bufio.NewReader 给它加上一层 4KB 的缓冲区，将多次小读取合并为少数系统调用，并提供 ReadString('\n') 这样便捷的按行读取方法——这对解析 SSE 这种逐行协议至关重要。这不是重复缓冲，而是必要的优化和功能增强。

body.Read(buf)，适合读取大块二进制的场景

# 四.SSE 常见问题
## 1. 忘记 Flush
错误：

```go
fmt.Fprintf(w, "data: hello\n\n")
```

正确：

```go
fmt.Fprintf(w, "data: hello\n\n")
flusher.Flush()
```

否则数据会停留在缓冲区。

---

## 2. Nginx 缓冲导致不流式
本地正常：

```latex
实时输出
```

线上变成：

```latex
一次性返回
```

通常是 Nginx 缓冲导致。

关闭缓冲：

```nginx
location /api/chat {

    proxy_buffering off;

    proxy_cache off;

    chunked_transfer_encoding on;
}
```



## 3. 用户中断连接
浏览器关闭页面后：

```go
select {
case <-ctx.Done():
    return
default:
}
```

及时终止模型调用，避免浪费 Token。

# 

