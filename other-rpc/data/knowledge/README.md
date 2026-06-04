# Go 语言 AI 助手知识库

本知识库面向 **Go 语言学习与开发** 场景的 AI 助手，覆盖版本信息、基础类型底层、并发、内存与常见实践。

---

## 助手定位

- 解答 Go 语法、标准库、工具链（`go`、`gofmt`、`go test`、`go mod`）相关问题。
- 说明值类型与引用类型的**底层表示**（栈/堆、指针、slice 头、map 桶等）。
- 对「当前最新版本」类问题，优先依据本知识库；若文档未更新，提示用户访问 https://go.dev/dl/ 核对。

---

## Go 当前版本（请以官网为准）

截至本知识库维护时：
当前go的版本已经发布到了1.26.3了，在2026.5.24

| 项目 | 说明 |
|------|------|
| 发布节奏 | 每年 2 月、8 月各一个 feature 版本 |
| 模块路径 | Go 1.11+ 默认使用 Go modules（`go.mod`） |
| 工具链命令 | `go version` 可查看本机安装版本 |

历史版本要点：Go 1.18 引入泛型；Go 1.20 起支持 WASI；Go 1.21 调整 `min`/`max` 等内置函数；Go 1.22 起 `for` 循环变量语义变更（每轮独立变量）。

---

## 基础类型与底层表示

### 整数与浮点

- **有符号整数**：`int8` `int16` `int32` `int64`；**无符号**：`uint8`（即 `byte`）`uint16` `uint32` `uint64`。
- **`int` 与 `uint`**：字长与平台相关（32 位机 32 位，64 位机 64 位）。
- **浮点**：`float32`、`float64`；运算中较小类型会提升为 `float64`。
- 底层为二进制补码（整数）与 IEEE 754（浮点），由编译器直接映射到机器字。

### bool、string

- `bool`：1 字节，值为 0 或 1（仅应使用 `true`/`false`）。
- **string**：只读字节序列；底层为 **指针 + 长度**（不暴露，语义上不可变）。与 `[]byte` 可零拷贝转换需注意：`string(b)` 会复制，`unsafe` 或 Go 1.20+ `unsafe.String` 等高级用法需谨慎。

### 数组与 slice

- **数组** `[N]T`：值类型，长度是类型的一部分，赋值/传参会复制整个数组。
- **slice** `[]T`：运行时结构（slice header）大致包含：
    - 指向底层数组的指针
    - `len`（长度）
    - `cap`（容量）
- `append` 在 `len == cap` 时可能触发扩容（通常按约 2 倍增长，具体由运行时实现）。
- **传 slice** 只复制 header，底层数组共享；修改元素会影响同一底层数组的其他 slice。

### map

- 引用类型，底层为 **哈希表**（桶 + 溢出桶，实现细节随版本演进）。
- **未初始化的 map**（`var m map[K]V`）不能写入，会 panic；须 `make(map[K]V)` 或字面量初始化。
- 遍历顺序**随机**，不要依赖顺序。
- key 必须可比较（不可含 slice、map、func 等）。

### struct 与指针

- 字段按声明顺序排列，存在 **内存对齐**；较大对齐要求的字段可能产生 padding。
- 指针 `*T` 存地址；`new(T)`、`&v` 获取指针。逃逸分析决定变量分配在栈还是堆。

### interface

- 底层为 **(type, value)** 二元组；`nil` interface 需 type 与 value 均为 nil 才为真 nil。
- 空接口 `interface{}` 或 `any`（Go 1.18+）可承载任意具体类型。
- 动态分发通过 **itable**（接口表）实现。

---

## 函数、方法与 defer

- 多返回值是常态；命名返回值在函数开头作用域内有效。
- `defer`：函数返回前 LIFO 执行，常用于关闭资源、`Unlock`。
- **方法**：值接收者与指针接收者影响是否复制 struct；大 struct 或需修改接收者时用指针接收者。

---

## 并发基础

### goroutine

- 轻量级用户态线程，由 Go 运行时调度（M:N 模型）。
- 启动：`go f()`；不要在没有同步的情况下共享可变状态。

### channel

- `chan T`：有缓冲 `make(chan T, n)` 与无缓冲之分。
- 无缓冲 channel：发送与接收必须同时就绪（同步握手）。
- 关闭：`close(ch)`；接收方可用 `v, ok := <-ch` 判断是否已关闭。
- 原则：**用 channel 传递所有权，用 mutex 保护共享状态**（按场景选择）。

### 常见同步

- `sync.Mutex` / `RWMutex`：互斥锁。
- `sync.WaitGroup`：等待一组 goroutine 结束。
- `context.Context`：取消、超时、传值（仅建议传请求作用域元数据）。

---

## 内存与 GC（简要）

- **垃圾回收**：并发标记清除，STW 时间已大幅缩短（具体因版本而异）。
- 避免在热路径频繁分配；可用 `sync.Pool` 复用对象。
- **逃逸分析**：编译器决定变量是否在堆上分配；可用 `go build -gcflags="-m"` 观察（学习用）。

---

## 错误处理

- 无 try/catch；使用 **`error` 接口**：`type error interface { Error() string }`。
- 惯用法：`if err != nil { return ..., err }`。
- Go 1.13+：`errors.Is`、`errors.As`；`fmt.Errorf("... %w", err)` 包装错误。
- panic 仅用于不可恢复或编程错误；库代码应优先返回 error。

---

## 模块与项目布局

- `go mod init example.com/foo` 创建模块。
- 依赖写在 `go.mod`；`go get`、`go mod tidy` 管理版本。
- 常见目录：`cmd/` 放 main，`internal/` 放私有包，`pkg/` 放可被外部引用的库（可选）。

---

## 性能与工具

- **pprof**：CPU、内存、阻塞分析。
- **基准测试**：`func BenchmarkXxx(b *testing.B)`。
- **竞态检测**：`go test -race`。
- **静态分析**：`go vet`、staticcheck、golangci-lint 等。

---

## 常见面试/实战要点

1. `slice` 与 `数组` 区别；`append` 扩容行为。
2. `map` 并发读写需 `sync.Map` 或加锁，否则 panic。
3. `interface` -nil 判断陷阱。
4. `select` 多路复用 channel；`default` 实现非阻塞。
5. **happens-before**：channel 操作、锁、`Once` 等建立同步关系。
6. **值拷贝 vs 指针**：slice/map/channel 本身是小型描述符，传参复制描述符但共享底层。

---
