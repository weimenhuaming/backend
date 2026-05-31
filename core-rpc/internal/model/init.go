package model

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(dsn string, logLevel logger.LogLevel, maxIdleConn, maxOpenConn int) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Printf("gorm open failed: %v", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("gorm get sql db failed: %v", err)
		os.Exit(1)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("mysql ping failed: %v", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	//if err := db.AutoMigrate(
	//	&entity.User{},
	//	&entity.Article{},
	//	&entity.Category{},
	//	&entity.Comment{},
	//	&entity.InteractionLike{},
	//	&entity.TokenBlacklist{},
	//); err != nil {
	//	log.Printf("auto migrate failed: %v", err)
	//	os.Exit(1)
	//}

	return db
}
