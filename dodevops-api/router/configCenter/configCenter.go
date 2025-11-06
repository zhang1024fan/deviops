// configCenter路由注册系统

package configCenter

import (
	"dodevops-api/api/configcenter/controller"

	"github.com/gin-gonic/gin"
)

// RegisterConfigCenterRoutes 注册配置中心相关路由
func RegisterConfigCenterRoutes(router *gin.RouterGroup) {
	ecsAuthCtrl := controller.NewEcsAuthController()
	accountAuthCtrl := controller.NewAccountAuthController()
	keyManageCtrl := controller.NewKeyManageController()
	syncScheduleCtrl := controller.NewSyncScheduleController()

	// ECS认证凭据管理
	router.GET("/config/ecsauthlist", ecsAuthCtrl.GetEcsAuthList)     // 获取所有凭据
	router.GET("/config/ecsauthinfo", ecsAuthCtrl.GetEcsAuthByName)   // 根据名称查找凭据
	router.GET("/config/ecsauthdetail", ecsAuthCtrl.GetEcsAuthById)   // 根据ID查找凭据详情
	router.POST("/config/ecsauthadd", ecsAuthCtrl.CreateEcsAuth)      // 创建凭据
	router.PUT("/config/ecsauthupdate", ecsAuthCtrl.UpdateEcsAuth)    // 更新凭据
	router.DELETE("/config/ecsauthdelete", ecsAuthCtrl.DeleteEcsAuth) // 删除凭据
	// 账号认证管理
	router.POST("/config/accountauth", accountAuthCtrl.Create)   // 创建账号
	router.PUT("/config/accountauth", accountAuthCtrl.Update)   // 更新账号
	router.DELETE("/config/accountauth", accountAuthCtrl.Delete)  // 删除账号
	router.GET("/config/accountauth", accountAuthCtrl.GetByID)    // 根据ID查询账号	
	router.GET("/config/accountauth/list", accountAuthCtrl.List)  // 获取所有账号
	router.POST("/config/accountauth/decrypt", accountAuthCtrl.DecryptPassword) // 解密密码
	router.GET("/config/accountauth/type", accountAuthCtrl.GetByType)    // 根据类型查询账号
	router.GET("/config/accountauth/alias", accountAuthCtrl.GetByAlias) // 根据别名查询账号

	// 密钥管理
	router.POST("/config/keymanage", keyManageCtrl.Create)              // 创建密钥
	router.PUT("/config/keymanage", keyManageCtrl.Update)               // 更新密钥
	router.DELETE("/config/keymanage", keyManageCtrl.Delete)            // 删除密钥
	router.GET("/config/keymanage", keyManageCtrl.GetByID)              // 根据ID查询密钥
	router.GET("/config/keymanage/list", keyManageCtrl.List)            // 获取密钥列表
	router.POST("/config/keymanage/decrypt", keyManageCtrl.DecryptKeys) // 解密密钥信息
	router.GET("/config/keymanage/type", keyManageCtrl.GetByType)       // 根据云厂商类型查询密钥
	router.POST("/config/keymanage/sync", keyManageCtrl.SyncCloudHosts) // 同步云主机（统一接口）

	// 定时同步配置管理
	router.POST("/config/sync-schedule", syncScheduleCtrl.Create)                      // 创建定时同步配置
	router.PUT("/config/sync-schedule", syncScheduleCtrl.Update)                       // 更新定时同步配置
	router.DELETE("/config/sync-schedule", syncScheduleCtrl.Delete)                    // 删除定时同步配置
	router.GET("/config/sync-schedule", syncScheduleCtrl.GetByID)                      // 根据ID查询定时同步配置
	router.GET("/config/sync-schedule/list", syncScheduleCtrl.List)                    // 获取定时同步配置列表
	router.POST("/config/sync-schedule/toggle-status", syncScheduleCtrl.ToggleStatus)  // 切换配置状态
	router.GET("/config/sync-schedule/active", syncScheduleCtrl.GetActiveSchedules)    // 获取启用的配置
	router.POST("/config/sync-schedule/trigger", syncScheduleCtrl.TriggerManualSync)   // 手动触发同步（测试用）
	router.GET("/config/sync-schedule/scheduler-stats", syncScheduleCtrl.GetSchedulerStats) // 获取调度器状态
	router.GET("/config/sync-schedule/log", syncScheduleCtrl.GetSyncLog)              // 获取同步日志
}
