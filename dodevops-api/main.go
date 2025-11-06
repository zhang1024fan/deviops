// 启动程序
// author xiaoRui

package main

import (
	"context"
	"flag"
	"fmt"
	"dodevops-api/common"
	"dodevops-api/common/config"
	_ "dodevops-api/docs"
	"dodevops-api/api/cmdb/controller"
	"dodevops-api/api/task/service"
	"dodevops-api/pkg/db"
	"dodevops-api/pkg/log"
	"dodevops-api/pkg/redis"
	"dodevops-api/router"
	"dodevops-api/scheduler"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

// @title devops运维管理系统
// @version 1.0
// @description devops运维管理系统API接口文档
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "c", "", "配置文件路径 (默认: ./config.yaml)")
	flag.Parse()

	// 加载配置文件
	if err := config.LoadConfig(configPath); err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化其他服务
	if err := initServices(); err != nil {
		panic("Failed to initialize services: " + err.Error())
	}

	// 启动调度器管理器
	if err := scheduler.GetManager().Start(); err != nil {
		log.Log().Errorf("Failed to start scheduler manager: %v", err)
	}

	// 初始化任务队列系统
	if err := initTaskQueue(); err != nil {
		log.Log().Errorf("Failed to initialize task queue: %v", err)
	}

	log := log.Log()
	gin.SetMode(config.Config.Server.Model)
	router := router.InitRouter()
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: router,
	}
	// 启动服务
	go func() {
		log.Info("Conflicting values for 'process.env.NODE_ENV'")
		log.Info("")
		log.Info("  App running at:")
		log.Info("  - Local:   http://" + config.Config.Server.Address)
		log.Info("  - Network: http://" + config.Config.Server.Address)
		log.Info("")
		log.Info("  请注意，开发版本尚未优化")
		log.Info("  要创建生产环境构建，请运行 go run main.go")
		log.Info("")
		if config.Config.Server.EnableSwagger {
			log.Info("API文档地址: http://" + config.Config.Server.Address + "/swagger/index.html")
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Infof("listen: %s \n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	//监听消息
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	// 停止调度器管理器
	scheduler.GetManager().Stop()

	// 停止任务队列系统
	stopTaskQueue()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Server Shutdown: %v", err)
	}
	log.Info("Server exiting")
}

// initServices 初始化各种服务
func initServices() error {
	// 初始化数据库连接和自动迁移
	if err := db.SetupDBLink(); err != nil {
		return err
	}

	// 初始化Redis
	if err := common.InitRedis(); err != nil {
		log.Log().Warnf("Redis initialization failed: %v", err)
	}
	redis.SetupRedisDb()

	// 初始化SQL记录控制器
	controller.InitCmdbSQLRecordController(common.GetDB())

	return nil
}

// initTaskQueue 初始化任务队列系统
func initTaskQueue() error {
	// 导入任务服务包
	taskService := service.NewTaskWorkService()
	_ = taskService // 避免未使用警告

	// 初始化全局调度器
	if err := service.InitGlobalScheduler(); err != nil {
		return fmt.Errorf("初始化全局调度器失败: %v", err)
	}

	// 初始化任务队列
	config := service.DefaultConfig()
	if err := service.InitTaskQueue(config); err != nil {
		return fmt.Errorf("初始化任务队列失败: %v", err)
	}

	log.Log().Info("任务队列系统初始化成功")
	return nil
}

// stopTaskQueue 停止任务队列系统
func stopTaskQueue() {
	log.Log().Info("停止任务队列系统...")

	// 停止全局调度器
	service.StopGlobalScheduler()

	// 停止任务队列
	service.StopTaskQueue()

	log.Log().Info("任务队列系统已停止")
}
