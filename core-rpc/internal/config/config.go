package config

import (
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm/logger"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql DatabaseConf
	Cache []redis.RedisConf
}

type DatabaseConf struct {
	Host        string `json:"host" yaml:"Host"`
	Port        string `json:"port" yaml:"Port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	DBName      string `json:"dbname" yaml:"dbname"`
	Config      string `json:"config" yaml:"config"`
	MaxIdleConn int    `json:"max_idle_conn" yaml:"Max_Idle_Conn"`
	MaxOpenConn int    `json:"max_open_conn" yaml:"Max_Open_Conn"`
	LogMode     string `json:"log_mode" yaml:"Log_Mode"`
}

func (m *DatabaseConf) Dsn() string {
	//dsn := "root:ccebdcxy@tcp(127.0.0.1:3306)/shorturl?charset=utf8&parseTime=True&loc=Local"
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DBName + "?" + m.Config
}

func (m DatabaseConf) LogLevel() logger.LogLevel {
	switch strings.ToLower(m.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
