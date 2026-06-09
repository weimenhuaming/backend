package svc

import (
	"core-rpc/core_client"
	"gateway/internal/config"
	"gateway/internal/middleware"
	"log"
	agent_client "other-rpc/agent_client"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Cache  *redis.Redis
	core_client.Core
	Agent          agent_client.Agent
	AuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	cache := redis.MustNewRedis(c.Cache[0])
	log.Println("Redis连接成功")

	return &ServiceContext{
		Config:         c,
		Core:           core_client.NewCore(zrpc.MustNewClient(c.CoreRpc)),
		Agent:          agent_client.NewAgent(zrpc.MustNewClient(c.AgentRpc)),
		Cache:          cache,
		AuthMiddleware: middleware.NewAuthMiddleware(c, cache).Handle,
	}
}
