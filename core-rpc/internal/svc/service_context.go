package svc

import (
	"core-rpc/internal/config"
	"core-rpc/internal/model"
	"core-rpc/internal/model/repo"
	"log"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	UserRepo *repo.UserModel
	Db       *gorm.DB
	Cache    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitGorm(c.Mysql.Dsn(), c.Mysql.LogLevel(), c.Mysql.MaxIdleConn, c.Mysql.MaxOpenConn)
	log.Println("gorm 连接成功，表结构已 AutoMigrate")

	cache := redis.MustNewRedis(c.Cache[0])
	log.Println("Redis 连接成功")

	return &ServiceContext{
		Config:   c,
		Db:       db,
		UserRepo: repo.NewUserModel(db),
		Cache:    cache,
	}
}
