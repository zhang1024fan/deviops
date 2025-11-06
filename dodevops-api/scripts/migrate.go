package main

import (
	"fmt"
	"dodevops-api/common"
	"dodevops-api/common/config"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/log"
)

func main() {
	// 初始化配置
	config.Setup()

	// 初始化日志
	log.Setup()

	// 初始化数据库连接
	if err := db.SetupDBLink(); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		panic(err)
	}

	// 执行数据库迁移
	if err := db.AutoMigrate(common.GetDB()); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		panic(err)
	}

	fmt.Println("数据库迁移成功完成！")
	fmt.Println("包含以下表：")
	fmt.Println("- applications")
	fmt.Println("- jenkins_envs")
	fmt.Println("- quick_deployments")
	fmt.Println("- quick_deployment_tasks")
	fmt.Println("以及其他系统表...")
}