// router/tool/tool.go
package tool

import (
	"dodevops-api/api/tool/controller"
	"github.com/gin-gonic/gin"
)

// RegisterToolRoutes 注册导航工具相关路由
func RegisterToolRoutes(router *gin.RouterGroup) {
	// 导航工具管理
	router.POST("/tool", controller.CreateTool)        // 创建导航工具
	router.GET("/tool/:id", controller.GetToolByID)    // 获取导航工具详情
	router.PUT("/tool", controller.UpdateTool)         // 更新导航工具
	router.DELETE("/tool/:id", controller.DeleteTool)  // 删除导航工具
	router.GET("/tool/list", controller.GetToolList)   // 获取导航工具列表（分页）
	router.GET("/tool/all", controller.GetAllTools)    // 获取所有导航工具（不分页）

	// 服务部署管理
	router.GET("/tool/services", controller.GetServicesList)             // 获取可部署服务列表
	router.GET("/tool/services/:serviceId", controller.GetServiceDetail) // 获取服务详情
	router.POST("/tool/deploy", controller.CreateDeploy)                 // 创建部署任务
	router.GET("/tool/deploy/list", controller.GetDeployList)            // 获取部署历史列表
	router.GET("/tool/deploy/:id/status", controller.GetDeployStatus)    // 获取部署状态
	router.DELETE("/tool/deploy/:id", controller.DeleteDeploy)           // 卸载服务
}
