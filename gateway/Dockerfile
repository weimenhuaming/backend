FROM golang:1.25.8-alpine3.23 AS builder

WORKDIR /build
COPY core-rpc/ ./core-rpc/
COPY gateway/go.mod gateway/go.sum ./gateway/

WORKDIR /build/gateway
RUN go mod download

COPY gateway/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/gateway gateway.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /build/gateway/etc/gateway-docker.yaml etc/gateway.yaml

EXPOSE 9000
CMD ["./gateway", "-f", "etc/gateway.yaml"]