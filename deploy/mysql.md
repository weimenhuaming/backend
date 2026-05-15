# Mysql的镜像下载

数据库并不会像es那些需要频繁安装组件，所以就直接使用docker run命令来下载
当然我这边依旧是挂载的一个设置，可以自行取消

```bash
docker run -d \
--name mysql8.4.7 \
-p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=123456 \
-e TZ=Asia/Shanghai \
-v "D:\Docker\mysql8.4.7:/var/lib/mysql" \
mysql:8.4.7
```