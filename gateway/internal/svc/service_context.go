package svc

import (
	"core-rpc/core"
	"gateway/internal/config"
	"log"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	core.Core              // 加入user rpc的服务操作函数
	Cache     *redis.Redis //使用gozero自带的，他本身就是对go-redis的一个封装
}

func NewServiceContext(c config.Config) *ServiceContext {
	cache := redis.MustNewRedis(c.Cache[0]) // 支持多结点，只用第一个
	log.Println("Redis连接成功")

	return &ServiceContext{
		Config: c,
		Core:   core.NewCore(zrpc.MustNewClient(c.CoreRpc)),
		Cache:  cache,
	}
}
