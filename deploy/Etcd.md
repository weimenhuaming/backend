# Etcd 镜像下载

## 一、Mac 上的下载

这里选择的是 etcd 官方镜像源： quay.io/coreos/etcd:v3.6.1

启动命令：

```bash
docker run -d \
  --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  -v $(pwd)/etcd-data:/etcd-data \
  quay.io/coreos/etcd:v3.6.1 \
  /usr/local/bin/etcd \
  --name etcd0 \
  --data-dir /etcd-data \
  --advertise-client-urls http://localhost:2379 \
  --listen-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380
```

一般来说在本地的测试环境，可以不要-v

-d：后台运行容器  
--name etcd：指定 Docker 容器名称  
-p：端口映射，将宿主机端口映射到容器内部端口  
-v：数据卷挂载，用于数据持久化  
/usr/local/bin/etcd：etcd 启动程序  
--name etcd0：指定当前 etcd 节点名称  
--data-dir：指定 etcd 数据存储目录  
--advertise-client-urls：指定 etcd 对外暴露给客户端访问的地址  
--listen-client-urls：指定 etcd 监听客户端请求的地址  
--listen-peer-urls：指定 etcd 监听集群节点通信的地址  
