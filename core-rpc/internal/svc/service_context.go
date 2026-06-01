package svc

import (
	"core-rpc/internal/config"
	"core-rpc/internal/model"
	"log"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Cache  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitGorm(c.Mysql.Dsn(), c.Mysql.LogLevel(), c.Mysql.MaxIdleConn, c.Mysql.MaxOpenConn)
	log.Println("gorm 连接成功，表结构已 AutoMigrate")

	cache := redis.MustNewRedis(c.Cache[0])
	log.Println("Redis 连接成功")

	return &ServiceContext{
		Config: c,
		Db:     db,
		Cache:  cache,
	}
}
