# v1.Elasticsearch + Kibana + IK分词器 的容器构建和启动记录

## 1. 启动容器
在 `docker-compose.yml` 所在目录执行：
``` bash
docker compose up -d
```

如果启动成功，可以使用以下命令查看：
``` bash
docker ps
```

## 2. 访问服务

### Elasticsearch
浏览器访问：http://localhost:9200 如果有json数据，则说明成功

### Kibana
浏览器访问：http://localhost:5601 Kibana 启动通常需要 **30 秒左右**。

## 3. 创建数据目录
为了防止容器删除后数据丢失，需要提前创建挂载目录，目录的位置你可以自行设置。

``` bash
mkdir elasticsearch-data
mkdir elasticsearch-plugin
```
说一下这两个目录的作用：  
elasticsearch-data     存储 Elasticsearch 索引数据  
elasticsearch-plugin   存储 Elasticsearch 插件  

## 4. 安装 IK 中文分词器

默认情况下，Elasticsearch 使用的是 **英文分词器**，中文分词效果非常差，因此需要安装 **IK Analyzer**。  
IK 是目前最常用的中文分词插件。
目前我测试下来这个容器里面是没有wget甚至apt这些命令，所以我们可以自行去官网下载压缩包，解压到elasticsearch-plugin  

ik分词器官网： https://github.com/infinilabs/analysis-ik 

解压之后的目录结构应该类似（如果没有ik，需要我们自行手建）：

    plugins
    └── ik
        ├── plugin-descriptor.properties
        ├── config
        └── ...


在安装完之后记得重启一下， 可以访问：http://localhost:9200/_cat/plugins
如果看到： elasticsearch analysis-ik 说明安装成功。

## 5. 测试 IK 分词
我们可以进入Kibana中，下划找到Dev Tools，输入下面json
``` json
{
  "analyzer": "ik_max_word",
  "text": "中华人民共和国"
}
```

返回示例：

    中华
    中华人民
    中华人民共和国
    人民
    共和国

说明 IK 分词器工作正常。


## 6. 可能遇到的问题

+ 首先在ES8.0之后，其实是默认打开安全验证的，由于我这边是测试环境，就没有开启，是需要ca证书的
+ 其它问题待处理

## 7. 附上代码

```bash
version: "3.7"
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:9.3.3
    container_name: elasticsearch
    restart: no
    environment:
      - xpack.security.enabled=false
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
    ulimits:
      memlock:
        soft: -1
        hard: -1
      nofile:
        soft: 65536
        hard: 65536
    cap_add:
      - IPC_LOCK
    volumes:
      - ./elasticsearch-data:/usr/share/elasticsearch/data
      - ./elasticsearch-plugin:/usr/share/elasticsearch/plugins
    ports:
      - 9200:9200

  kibana:
    container_name: kibana
    image: docker.elastic.co/kibana/kibana:9.3.3
    restart: no
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - 5601:5601
    depends_on:
      - elasticsearch
```

# v2. 更为标准的做法
其实上述v1的版本并不是特别好，因为每一次都需要我们自行下载组件，当然应该也可以在docker-compose里直接下载，我查完资料感觉这样写有点难看。  
我推荐还是将es先构建，加入完分词之后，在通过docker-compose完成这个整体的搭建的一个过程吧

## 1. 搭建elasticsearch

我这边已经搭建好了，并且也下载了ik分词器，唯一的问题就是如果把组件plugin的部分进行挂载，就会导致内容覆盖的问题，说简单就是之前在镜像里下载好的ik，会因为挂载覆盖点，从而丢失。  
所以我们就不挂载plugin了，如果需要，我们后续再去修改es的dockerfile

## 2. docker-compose的修改

修改好的我已经放在当前目录下了，自行查看区别


## 3. 构建

这个时候的构建命令就需要做一下更改了：
```bash
    docker compose up -d --build
```


# 三.说明和问题

当前配置主要用于：

-   本地开发
-   测试环境
-   学习 Elasticsearch

配置特点：

-   单节点模式
-   关闭安全认证
-   挂载数据目录
-   支持中文 IK 分词器

生产环境一般需要：

-   开启安全认证
-   多节点集群
-   专用存储
-   监控与备份
