package model

import (
	"os"
	"time"
)

func InitGorm(Dsn string, LogLevel logger.LogLevel, MaxIdleConn, MaxOpenConn int) *gorm.DB {
	db, err := gorm.Open(mysql.Open(Dsn), &gorm.Config{
		Logger: logger.Default.LogMode(LogLevel),
	})
	//db.AutoMigrate(&LeaveMsg{})

	if err != nil {
		os.Exit(1)
	}

	// 获取底层的 SQL 数据库连接对象
	sqlDB, _ := db.DB()

	// 看看能不能ping的通
	if err := sqlDB.Ping(); err != nil {
		os.Exit(1)
	}

	// 设置数据库连接池中的最大空闲连接数
	sqlDB.SetMaxIdleConns(MaxIdleConn) //10
	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(MaxOpenConn) //100
	//sqlDB.SetConnMaxLifetime(time.Hour * 7)     // 强制回收
	sqlDB.SetConnMaxIdleTime(time.Minute * 30) // 空闲连接可以存在的最长时间
	return db
}
