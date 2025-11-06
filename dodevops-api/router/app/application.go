package app

import (
	"dodevops-api/api/app/controller"
	"dodevops-api/common"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterApplicationRoutes 注册应用管理路由
func RegisterApplicationRoutes(router *gin.RouterGroup) {
	appCtrl := controller.NewApplicationController(common.GetDB())

	// 应用管理路由
	router.POST("/apps", middleware.AuthMiddleware(), appCtrl.CreateApplication)                    // 创建应用
	router.GET("/apps", middleware.AuthMiddleware(), appCtrl.GetApplicationList)                   // 获取应用列表
	router.GET("/apps/:id", middleware.AuthMiddleware(), appCtrl.GetApplicationDetail)             // 获取应用详情
	router.PUT("/apps/:id", middleware.AuthMiddleware(), appCtrl.UpdateApplication)                // 更新应用
	router.DELETE("/apps/:id", middleware.AuthMiddleware(), appCtrl.DeleteApplication)             // 删除应用

	// 应用的Jenkins环境配置管理路由（与应用强绑定）
	router.GET("/apps/:id/jenkins-envs", middleware.AuthMiddleware(), appCtrl.GetAppJenkinsEnvs)        // 获取应用的Jenkins环境配置
	router.POST("/apps/:id/jenkins-envs", middleware.AuthMiddleware(), appCtrl.AddAppJenkinsEnv)        // 添加应用的Jenkins环境配置
	router.PUT("/apps/:id/jenkins-envs/:env_id", middleware.AuthMiddleware(), appCtrl.UpdateAppJenkinsEnv)    // 更新应用的Jenkins环境配置
	router.DELETE("/apps/:id/jenkins-envs/:env_id", middleware.AuthMiddleware(), appCtrl.DeleteAppJenkinsEnv) // 删除应用的Jenkins环境配置

	// Jenkins服务器相关路由
	router.GET("/apps/jenkins-servers", middleware.AuthMiddleware(), appCtrl.GetJenkinsServers) // 获取Jenkins服务器列表
	router.POST("/apps/jenkins-job/validate", middleware.AuthMiddleware(), appCtrl.ValidateJenkinsJob) // 验证Jenkins任务是否存在

	// 业务线服务树路由
	router.GET("/apps/service-tree", middleware.AuthMiddleware(), appCtrl.GetServiceTree) // 获取业务线服务树
	router.GET("/apps/business-group-options", middleware.AuthMiddleware(), appCtrl.GetBusinessGroupOptions) // 获取业务组选项（连级选择器）
	router.GET("/apps/environment", middleware.AuthMiddleware(), appCtrl.GetAppEnvironment) // 获取单个应用环境配置


	// 快速发布路由
	router.GET("/apps/deployment/applications", middleware.AuthMiddleware(), appCtrl.GetApplicationsForDeployment) // 获取可发布应用列表
	router.POST("/apps/deployment/quick", middleware.AuthMiddleware(), appCtrl.CreateQuickDeployment)              // 创建快速发布
	router.POST("/apps/deployment/execute", middleware.AuthMiddleware(), appCtrl.ExecuteQuickDeployment)           // 执行快速发布
	router.GET("/apps/deployment/list", middleware.AuthMiddleware(), appCtrl.GetQuickDeploymentList)               // 获取快速发布列表
	router.GET("/apps/deployment/:id", middleware.AuthMiddleware(), appCtrl.GetQuickDeploymentDetail)              // 获取快速发布详情
	router.DELETE("/apps/deployment/:id", middleware.AuthMiddleware(), appCtrl.DeleteQuickDeployment)              // 删除快速发布（级联删除子任务）
	router.GET("/apps/deployment/tasks/:task_id/log", middleware.AuthMiddleware(), appCtrl.GetTaskBuildLog)        // 获取任务构建日志
	router.GET("/apps/deployment/tasks/:task_id/status", middleware.AuthMiddleware(), appCtrl.GetTaskStatus)       // 获取任务状态
}

// RegisterJenkinsRoutes 注册Jenkins路由
func RegisterJenkinsRoutes(router *gin.RouterGroup) {
	jenkinsCtrl := controller.NewJenkinsController(common.GetDB())

	// Jenkins服务器管理路由
	router.GET("/jenkins/servers", middleware.AuthMiddleware(), jenkinsCtrl.GetJenkinsServers)                    // 获取Jenkins服务器列表
	router.GET("/jenkins/servers/:id", middleware.AuthMiddleware(), jenkinsCtrl.GetJenkinsServerDetail)          // 获取Jenkins服务器详情
	router.POST("/jenkins/test-connection", middleware.AuthMiddleware(), jenkinsCtrl.TestJenkinsConnection)      // 测试Jenkins连接

	// Jenkins任务管理路由 (使用不同的路径结构)
	router.GET("/jenkins/:serverId/jobs", middleware.AuthMiddleware(), jenkinsCtrl.GetJobs)                       // 获取任务列表
	router.GET("/jenkins/:serverId/jobs/search", middleware.AuthMiddleware(), jenkinsCtrl.SearchJobs)             // 搜索任务
	router.GET("/jenkins/:serverId/jobs/:jobName", middleware.AuthMiddleware(), jenkinsCtrl.GetJobDetail)         // 获取任务详情

	// Jenkins任务操作路由
	router.POST("/jenkins/:serverId/jobs/:jobName/start", middleware.AuthMiddleware(), jenkinsCtrl.StartJob)      // 启动任务

	// Jenkins构建管理路由
	router.GET("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber", middleware.AuthMiddleware(), jenkinsCtrl.GetBuildDetail)               // 获取构建详情
	router.POST("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber/stop", middleware.AuthMiddleware(), jenkinsCtrl.StopBuild)             // 停止构建
	router.GET("/jenkins/:serverId/jobs/:jobName/builds/:buildNumber/log", middleware.AuthMiddleware(), jenkinsCtrl.GetBuildLog)             // 获取构建日志

	// Jenkins系统信息路由
	router.GET("/jenkins/:serverId/system-info", middleware.AuthMiddleware(), jenkinsCtrl.GetSystemInfo)         // 获取系统信息
	router.GET("/jenkins/:serverId/queue", middleware.AuthMiddleware(), jenkinsCtrl.GetQueueInfo)                // 获取队列信息
}