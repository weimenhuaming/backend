# 使用 Alpine Linux（更轻量，Windows 兼容性更好）
FROM alpine:latest
LABEL authors="chenaqi"

# 使用 ARG 定义构建参数（不会进入最终镜像的环境变量，避免与 etcd 冲突）
ARG ETCD_VERSION=3.6.10

# 设置 etcd 运行时需要的环境变量
ENV ETCD_DATA_DIR=/etcd-data

# 安装必要工具并下载 etcd
RUN apk add --no-cache wget tar \
    && wget https://github.com/etcd-io/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz \
    && tar -xzvf etcd-v${ETCD_VERSION}-linux-amd64.tar.gz \
    && mv etcd-v${ETCD_VERSION}-linux-amd64/etcd* /usr/local/bin/ \
    && rm -rf etcd-v${ETCD_VERSION}-linux-amd64* \
    && mkdir -p ${ETCD_DATA_DIR}

# 暴露端口
EXPOSE 2379 2380

# 声明数据卷
VOLUME /etcd-data

# 启动 etcd（移除 CMD 中的 --data-dir，因为已通过 ENV ETCD_DATA_DIR 定义）
ENTRYPOINT ["/usr/local/bin/etcd"]
CMD ["--advertise-client-urls", "http://0.0.0.0:2379", \
     "--listen-client-urls", "http://0.0.0.0:2379"]