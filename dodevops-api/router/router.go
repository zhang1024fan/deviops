package router
// utils/router.go

import (
	"dodevops-api/api/system/controller"
	"dodevops-api/common/config"
	"dodevops-api/middleware"
	"dodevops-api/pkg/log"

	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"dodevops-api/router/cmdb"         // cmdb模块路由
	"dodevops-api/router/configCenter" // 配置中心模块路由
	"dodevops-api/router/app"          // app模块路由
	"dodevops-api/router/dashboard"    // 看板模块路由
	"dodevops-api/router/k8s"          // k8s模块路由
	"dodevops-api/router/monitor"      // 监控模块路由
	"dodevops-api/router/system"       // 系统模块路由
	"dodevops-api/router/task"         // 任务中心路由
	"dodevops-api/router/tool"         // 导航工具路由

	agentController "dodevops-api/api/monitor/controller" // Agent控制器
)

// 初始化路由
func InitRouter() *gin.Engine {
	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 告诉 validator 使用 json 标签作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			tag := fld.Tag.Get("json")
			if tag == "-" {
				return ""
			}
			// 去掉 json 标签里的 ",omitempty"
			return strings.Split(tag, ",")[0]
		})
	}

	// 中间件
	router.Use(gin.Recovery())
	router.Use(middleware.Cors())
	// 让 /api/v1/upload/ 开头的请求映射到 ./upload/ 目录下
	// 使用绝对路径避免容器内相对路径问题
	uploadDir := config.Config.ImageSettings.UploadDir
	if !strings.HasPrefix(uploadDir, "/") {
		// 相对路径转绝对路径
		if absPath, err := filepath.Abs(uploadDir); err == nil {
			uploadDir = absPath
			log.Log().Infof("Static upload directory: %s -> %s", config.Config.ImageSettings.UploadDir, uploadDir)
		}
	} else {
		log.Log().Infof("Static upload directory: %s", uploadDir)
	}
	router.Static("/api/v1/upload", uploadDir)
	router.Use(log.CustomGinLogger())

	// 路由注册
	register(router)
	
	// 单独注册WebSocket路由到根路径
	registerWebSocketRoutes(router)

	return router
}

// 路由注册中心
func register(router *gin.Engine) {
	// 公共接口：Swagger、静态资源等
	// 根据配置决定是否启用Swagger文档
	if config.Config.Server.EnableSwagger {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Log().Info("Swagger API文档已启用")
	} else {
		log.Log().Info("Swagger API文档已禁用")
	}
	//router.Static("/upload", config.Config.ImageSettings.UploadDir)

	// 初始化控制器
	agentCtrl := agentController.NewAgentController()

	// 统一 API 分组
	apiGroup := router.Group("/api/v1")
	{
		// 不需要 JWT 的接口
		apiGroup.GET("/captcha", controller.Captcha)  // 验证码接口
		apiGroup.POST("/login", controller.Login)    // 登录接口
		// Agent心跳接口 - 不需要认证
		apiGroup.POST("/monitor/agent/heartbeat", agentCtrl.UpdateHeartbeat)
		
		// 需要 JWT鉴权 的接口
		jwtGroup := apiGroup.Group("")
		jwtGroup.Use(middleware.AuthMiddleware())
		jwtGroup.Use(middleware.LogMiddleware())
		{
			system.RegisterSystemRoutes(jwtGroup)
			cmdb.RegisterCmdbRoutes(jwtGroup)
			configCenter.RegisterConfigCenterRoutes(jwtGroup)
			app.RegisterJenkinsRoutes(jwtGroup)         // Jenkins模块路由
			app.RegisterApplicationRoutes(jwtGroup)     // 应用管理路由
			dashboard.RegisterDashboardRoutes(jwtGroup) // 看板模块路由
			k8s.RegisterK8sRoutes(jwtGroup)
			monitor.InitMonitorRouter(jwtGroup) // 新增监控路由
			task.RegisterTaskRoutes(jwtGroup)   // 任务中心路由
			tool.RegisterToolRoutes(jwtGroup)   // 导航工具路由
		}
	}
}

// registerWebSocketRoutes 注册WebSocket路由到根路径
func registerWebSocketRoutes(router *gin.Engine) {
	// WebSocket路由直接注册到根路径，前端可以访问 /ws/task/ansible/{id}/log/{work_id}
	task.RegisterWebSocketRoutes(router.Group(""))
}
