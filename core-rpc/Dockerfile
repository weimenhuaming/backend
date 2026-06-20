FROM docker.m.daocloud.io/library/golang:1.25.8-alpine3.23 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE=""

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/core-rpc core.go

FROM docker.m.daocloud.io/library/alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/core-rpc .
COPY etc/core-docker.yaml etc/core.yaml

EXPOSE 8080
CMD ["./core-rpc", "-f", "etc/core.yaml"]