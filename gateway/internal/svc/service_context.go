package svc

import (
	"core-rpc/core_client"
	"gateway/internal/config"
	"gateway/internal/middleware"
	"log"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	Cache            *redis.Redis //使用gozero自带的，他本身就是对go-redis的一个封装
	core_client.Core              // 加入user rpc的服务操作函数
	AuthMiddleware   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	cache := redis.MustNewRedis(c.Cache[0]) // 支持多结点，只用第一个
	log.Println("Redis连接成功")

	return &ServiceContext{
		Config:         c,
		Core:           core_client.NewCore(zrpc.MustNewClient(c.CoreRpc)),
		Cache:          cache,
		AuthMiddleware: middleware.NewAuthMiddleware(c).Handle, // 创建实例，这边相当于是服务端
	}
}
