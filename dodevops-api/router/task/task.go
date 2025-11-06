package task

import (
	"dodevops-api/api/task/controller"
	"dodevops-api/api/task/service"
	"dodevops-api/common"
	"dodevops-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.RouterGroup) {
	// 任务模板路由
	router.POST("/template/add", middleware.AuthMiddleware(), controller.CreateTemplate)
	router.GET("/template/list", middleware.AuthMiddleware(), controller.GetAllTemplates)
	router.PUT("/template/update", middleware.AuthMiddleware(), controller.UpdateTemplate)
	router.DELETE("/template/delete", middleware.AuthMiddleware(), controller.DeleteTemplate)
	router.GET("/template/info/:id", middleware.AuthMiddleware(), controller.GetTemplateByID)
	router.GET("/template/content/:id", middleware.AuthMiddleware(), controller.GetTemplateContent)
	router.GET("/template/query/name", middleware.AuthMiddleware(), controller.GetTemplatesByName)
	router.GET("/template/query/type", middleware.AuthMiddleware(), controller.GetTemplatesByType)

	// 任务管理路由
	router.POST("/task/add", middleware.AuthMiddleware(), controller.CreateTask)                      // 创建任务
	router.GET("/task/get", middleware.AuthMiddleware(), controller.GetTaskByID)                      // 获取任务信息
	router.PUT("/task/update", middleware.AuthMiddleware(), controller.UpdateTask)                    // 修改任务
	router.DELETE("/task/delete", middleware.AuthMiddleware(), controller.DeleteTask)                 // 删除任务
	router.GET("/task/list", middleware.AuthMiddleware(), controller.ListTasks)                       // 获取任务列表
	router.GET("/task/list-with-details", middleware.AuthMiddleware(), controller.ListTasksWithDetails) // 获取任务列表（包含关联信息）
	router.GET("/task/query/name", middleware.AuthMiddleware(), controller.GetTasksByName)            // 获取任务名称列表
	router.GET("/task/query/type", middleware.AuthMiddleware(), controller.GetTasksByType)            // 获取任务类型列表
	router.GET("/task/query/status", middleware.AuthMiddleware(), controller.GetTasksByStatus)         // 获取任务状态列表
	router.GET("/task/next-execution", middleware.AuthMiddleware(), controller.GetNextExecutionTime)  // 获取任务下次执行时间
	router.GET("/task/execution-info", middleware.AuthMiddleware(), controller.GetTaskExecutionInfo) // 获取任务执行信息
	router.GET("/task/templates", middleware.AuthMiddleware(), controller.GetTaskTemplatesWithStatus) // 获取任务模板列表

	// 任务作业路由
	router.POST("/taskjob/start", middleware.AuthMiddleware(), controller.TaskWork().StartJob)
	router.GET("/taskjob/log", middleware.AuthMiddleware(), controller.TaskWork().GetJobLog)
	router.POST("/taskjob/stop", middleware.AuthMiddleware(), controller.TaskWork().StopJob)
	router.GET("/taskjob/status", middleware.AuthMiddleware(), controller.TaskWork().GetJobStatus)

	// 任务监控路由
	taskMonitorCtrl := controller.NewTaskMonitorController()
	router.GET("/task/monitor/queue/metrics", middleware.AuthMiddleware(), taskMonitorCtrl.GetQueueMetrics)      // 获取队列指标
	router.GET("/task/monitor/scheduler/stats", middleware.AuthMiddleware(), taskMonitorCtrl.GetSchedulerStats) // 获取调度器统计
	router.GET("/task/monitor/system/status", middleware.AuthMiddleware(), taskMonitorCtrl.GetSystemStatus)     // 获取系统状态
	router.GET("/task/monitor/queue/details", middleware.AuthMiddleware(), taskMonitorCtrl.GetQueueDetails)     // 获取队列详情
	router.POST("/task/monitor/queue/clear-failed", middleware.AuthMiddleware(), taskMonitorCtrl.ClearFailedQueue) // 清空失败队列
	router.POST("/task/monitor/queue/retry-failed", middleware.AuthMiddleware(), taskMonitorCtrl.RetryFailedTasks) // 重试失败任务

	// 定时任务管理路由
	router.POST("/task/monitor/scheduled/pause", middleware.AuthMiddleware(), taskMonitorCtrl.PauseScheduledTask)  // 暂停定时任务
	router.POST("/task/monitor/scheduled/resume", middleware.AuthMiddleware(), taskMonitorCtrl.ResumeScheduledTask) // 恢复定时任务
	router.POST("/task/monitor/scheduled/reset", middleware.AuthMiddleware(), taskMonitorCtrl.ResetScheduledTaskStatus) // 重置定时任务状态
	router.GET("/task/monitor/task/status", middleware.AuthMiddleware(), taskMonitorCtrl.GetTaskStatus)           // 获取任务状态详情

	// Ansible任务路由
	taskAnsibleCtrl := controller.NewTaskAnsibleController(service.NewTaskAnsibleService(common.GetDB()))
	router.GET("/task/ansiblelist", middleware.AuthMiddleware(), taskAnsibleCtrl.List)                   // 获取任务列表
	router.POST("/task/ansible", middleware.AuthMiddleware(), taskAnsibleCtrl.CreateTask)                // 创建Ansible任务
	router.POST("/task/k8s", middleware.AuthMiddleware(), taskAnsibleCtrl.CreateK8sTask)                 // 创建K8s任务
	router.GET("/task/ansible/:id", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTask)                // 通过任务ID获取任务信息
	router.POST("/task/ansible/:id/start", middleware.AuthMiddleware(), taskAnsibleCtrl.StartTask)       // 启动任务(支持Ansible和K8s)
	router.DELETE("/task/ansible/:id", middleware.AuthMiddleware(), taskAnsibleCtrl.DeleteTask)          // 删除任务(级联删除子任务)
	router.GET("/task/ansible/:id/log/:work_id", middleware.AuthMiddleware(), taskAnsibleCtrl.GetJobLog) // 获取任务日志(SSE)
	router.GET("/task/ansible/query/name", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTasksByName)  // 根据名称模糊查询任务
	router.GET("/task/ansible/query/type", middleware.AuthMiddleware(), taskAnsibleCtrl.GetTasksByType)  // 根据类型查询任务

}

// RegisterWebSocketRoutes 注册WebSocket路由（不需要中间件认证）
func RegisterWebSocketRoutes(router *gin.RouterGroup) {
	// WebSocket任务日志路由 (内部处理认证)
	wsCtrl := controller.NewWebSocketController(service.NewTaskAnsibleService(common.GetDB()))
	router.GET("/ws/task/ansible/:id/log/:work_id", wsCtrl.GetJobLogWS) // WebSocket日志接口
}
