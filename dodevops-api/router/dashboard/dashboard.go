package dashboard

import (
	"dodevops-api/api/dashboard/controller"

	"github.com/gin-gonic/gin"
)

// RegisterDashboardRoutes 注册看板路由
func RegisterDashboardRoutes(router *gin.RouterGroup) {
	dashboardController := controller.NewDashboardController()

	// 看板路由组
	dashboardGroup := router.Group("/dashboard")
	{
		dashboardGroup.GET("/stats", dashboardController.GetDashboardStats) // 获取看板统计数据
		dashboardGroup.GET("/business-distribution", dashboardController.GetBusinessDistributionStats) // 获取业务分布统计数据
		dashboardGroup.GET("/assets", dashboardController.GetAssetStats) // 获取资产统计数据
	}
}