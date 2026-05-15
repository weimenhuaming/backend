# Redis 的镜像下载

Redis 和数据库类似，一般也不需要额外安装插件，所以直接使用 docker run 启动容器即可。
当然我这边依旧是挂载的一个设置，可以自行取消。

如果你想开启持久化，可以加上这个参数  --appendonly yes 表示开启 AOF 持久化。

```bash
docker run -d \
--name redis7 \
-p 6379:6379 \
-e TZ=Asia/Shanghai \
-v "D:\Docker\redis7:/data" \
redis:7
```