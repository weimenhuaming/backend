package svc

import (
	"core-rpc/internal/config"
	"core-rpc/internal/model/user"
	"log"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel user.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	log.Println("MySQL连接成功")
	return &ServiceContext{
		Config:    c,
		UserModel: user.NewUserModel(conn),
	}
}
