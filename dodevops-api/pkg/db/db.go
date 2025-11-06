// pkg/db/db.go
package db

import (
	"fmt"
	"dodevops-api/common/config"
	"io"
	"log"
	"os"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

// 创建自定义GORM日志记录器，同时输出到控制台和文件
func NewGormLogger() logger.Interface {
	// 确保logs目录存在
	os.MkdirAll("logs", os.ModePerm)

	// 打开日志文件
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("无法打开日志文件: %v", err)
		return logger.Default.LogMode(logger.Info)
	}

	// 创建多重写入器，同时写入控制台和文件
	mw := io.MultiWriter(os.Stdout, file)

	return logger.New(
		log.New(mw, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Info,   // 日志级别
			IgnoreRecordNotFoundError: false,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,          // 彩色打印
		},
	)
}

// 数据库初始化
func SetupDBLink() error {
	var err error
	var dbConfig = config.Config.Db
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&sql_mode=''",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Db,
		dbConfig.Charset)
	Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:                                   NewGormLogger(),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	if Db.Error != nil {
		panic(Db.Error)
	}

	// 自动建表
	if err := AutoMigrate(Db); err != nil {
		panic(err)
	}

	sqlDB, err := Db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)
	return nil
}
