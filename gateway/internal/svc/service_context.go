package svc

import (
	"core-rpc/core"
	"gateway/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	core.Core // 加入user rpc的服务操作函数
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Core:   core.NewCore(zrpc.MustNewClient(c.CoreRpc)),
	}
}
