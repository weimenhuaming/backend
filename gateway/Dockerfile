# 多阶段构建：编译阶段 + 运行阶段，最终镜像只保留二进制和必要文件
# 注意：docker-compose 的 build context 为仓库根目录（..），因此路径以根目录为基准

# ---------- 编译阶段 ----------
FROM golang:1.25.8-alpine3.23 AS builder

# Go 模块与国内代理配置，加速依赖下载
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE=""

WORKDIR /build

# gateway 通过 replace 引用本地 core-rpc / other-rpc，需一并复制才能解析依赖
COPY core-rpc/ ./core-rpc/
COPY other-rpc/ ./other-rpc/
COPY gateway/go.mod gateway/go.sum ./gateway/

WORKDIR /build/gateway
RUN go mod download

COPY gateway/ ./
# CGO_ENABLED=0 生成静态链接二进制；-ldflags="-s -w" 去除调试信息以减小体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/gateway gateway.go

# ---------- 运行阶段 ----------
FROM alpine:3.20

# CA 证书（HTTPS）与时区数据
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/gateway .
# 使用 Docker 专用配置（etcd/redis 等服务名与本地 gateway.yaml 不同）
COPY --from=builder /build/gateway/etc/gateway-docker.yaml etc/gateway.yaml
# 头像、博客图片等静态资源目录
RUN mkdir -p static/avatars static/blog

EXPOSE 9000
CMD ["./gateway", "-f", "etc/gateway.yaml"]
