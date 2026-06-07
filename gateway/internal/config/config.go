package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	CoreRpc  zrpc.RpcClientConf
	AgentRpc zrpc.RpcClientConf
	Auth     struct {
		AccessSecret string
		AccessExpire uint64
	}
	RefreshSecret string
	RefreshExpire uint64

	Cache []redis.RedisConf
}
