FROM golang:1.25.8-alpine3.23 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/other-rpc agent.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/other-rpc .
COPY etc/agent-docker.yaml etc/agent.yaml
COPY data/knowledge ./data/knowledge

EXPOSE 8081
CMD ["./other-rpc", "-f", "etc/agent.yaml"]