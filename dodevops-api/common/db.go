package common

import (
	"fmt"
	"dodevops-api/pkg/db"

	"gorm.io/gorm"
)

// GetDB 获取数据库连接（使用pkg/db中的连接）
func GetDB() *gorm.DB {
	if db.Db == nil {
		panic("Database connection is not initialized")
	}

	sqlDB, err := db.Db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database instance: %v", err))
	}

	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("Database connection lost: %v", err))
	}

	return db.Db
}
