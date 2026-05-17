package svc

import (
	"core-rpc/internal/config"
	"core-rpc/internal/model/article"
	"core-rpc/internal/model/user"
	"log"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	UserModel    user.UserModel
	ArticleModel article.ArticleModel
	Cache        *redis.Redis //使用gozero自带的，他本身就是对go-redis的一个封装
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	log.Println("MySQL连接成功")

	cache := redis.MustNewRedis(c.Cache[0]) // 支持多结点，只用第一个
	log.Println("Redis连接成功")
	return &ServiceContext{
		Config:       c,
		UserModel:    user.NewUserModel(conn),
		ArticleModel: article.NewArticleModel(conn),
		Cache:        cache,
	}
}
