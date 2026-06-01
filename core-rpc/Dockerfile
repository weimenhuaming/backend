FROM golang:1.25.8-alpine3.23 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/core-rpc core.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/core-rpc .
COPY etc/core-docker.yaml etc/core.yaml

EXPOSE 8080
CMD ["./core-rpc", "-f", "etc/core.yaml"]