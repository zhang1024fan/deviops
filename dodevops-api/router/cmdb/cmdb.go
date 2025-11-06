package cmdb

import (
	"dodevops-api/api/cmdb/controller"
	"dodevops-api/api/cmdb/service"

	"github.com/gin-gonic/gin"
)

// RegisterCmdbRoutes 注册系统相关路由
func RegisterCmdbRoutes(router *gin.RouterGroup) {
	// 资产分组
	router.POST("/cmdb/groupadd", controller.CreateCmdbGroup)                    // 添加资产分组
	router.GET("/cmdb/grouplist", controller.GetAllCmdbGroups)                   // 获取所有资产分组树
	router.GET("/cmdb/grouplistwithhosts", controller.GetAllCmdbGroupsWithHosts) // 获取所有资产分组树及关联主机
	router.PUT("/cmdb/groupupdate", controller.UpdateCmdbGroup)                  // 更新资产分组
	router.DELETE("/cmdb/groupdelete", controller.DeleteCmdbGroup)               // 删除资产分组
	router.GET("/cmdb/groupbyname", controller.GetCmdbGroupByName)               // 根据名称查询分组
	// 主机管理
	router.POST("/cmdb/hostcreate", controller.NewCmdbHostController().CreateCmdbHost)            // 创建主机
	router.PUT("/cmdb/hostupdate", controller.NewCmdbHostController().UpdateCmdbHost)             // 更新主机
	router.DELETE("/cmdb/hostdelete", controller.NewCmdbHostController().DeleteCmdbHost)          // 删除主机
	router.GET("/cmdb/hostlist", controller.NewCmdbHostController().GetCmdbHostListWithPage)      // 获取主机列表(分页)
	router.GET("/cmdb/hostinfo", controller.NewCmdbHostController().GetCmdbHostById)              // 根据ID获取主机
	router.GET("/cmdb/hostgroup", controller.NewCmdbHostController().GetCmdbHostsByGroupId)       // 根据分组ID获取主机列表
	router.GET("/cmdb/hostbyname", controller.NewCmdbHostController().GetCmdbHostsByHostNameLike) // 根据主机名称模糊查询
	router.GET("/cmdb/hostbyip", controller.NewCmdbHostController().GetCmdbHostsByIP)             // 根据IP查询主机
	router.GET("/cmdb/hostbystatus", controller.NewCmdbHostController().GetCmdbHostsByStatus)     // 根据状态查询主机
	router.POST("/cmdb/hostimport", controller.NewCmdbHostController().ImportHostsFromExcel)      // 从Excel导入主机
	router.GET("/cmdb/hosttemplate", controller.NewCmdbHostController().DownloadHostTemplate)     // 下载主机导入模板
	router.POST("/cmdb/hostsync", controller.NewCmdbHostController().SyncHostInfo)                // 同步主机基本信息
	// 云主机管理
	router.POST("/cmdb/hostcloudcreatealiyun", controller.NewCmdbHostCloudController().CreateAliyunHost)                          // 创建阿里云主机
	router.POST("/cmdb/hostcloudcreatetencent", controller.NewCmdbHostCloudController().CreateTencentHost)                        // 创建腾讯云主机
	router.POST("/cmdb/hostcloudcreatebaidu", controller.NewCmdbHostCloudController().CreateBaiduHost)                            // 创建百度云主机
	router.GET("/cmdb/hostssh/connect/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ConnectTerminal) // SSH终端连接
	router.GET("/cmdb/hostssh/command/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).ExecuteCommand)  // SSH执行命令
	router.POST("/cmdb/hostssh/upload/:id", controller.NewCmdbHostSSHController(service.GetCmdbHostSSHService()).UploadFile)      // SSH文件上传
	// SQL执行
	router.POST("/cmdb/sql/select", controller.GetCmdbSQLRecordController().ExecuteSelect)       // 执行查询语句(通过数据库ID/名称)
	router.POST("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteInsert)              // 执行插入语句(通过数据库ID/名称)
	router.PUT("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteUpdate)               // 执行更新语句(通过数据库ID/名称)
	router.DELETE("/cmdb/sql", controller.GetCmdbSQLRecordController().ExecuteDelete)            // 执行删除语句(通过数据库ID/名称)
	router.POST("/cmdb/sql/execute", controller.GetCmdbSQLRecordController().ExecuteSQL)         // 执行原生SQL语句(通过数据库ID/名称)
	router.POST("/cmdb/sql/databaselist", controller.GetCmdbSQLRecordController().ListDatabases) // 获取数据库列表(通过数据库ID)
	// SQL日志管理
	router.GET("/cmdb/sqlLog/list", controller.GetCmdbSqlLogList)         // 分页获取SQL操作日志列表
	router.DELETE("/cmdb/sqlLog/delete", controller.DeleteCmdbSqlLogById) // 根据id删除SQL操作日志
	router.DELETE("/cmdb/sqlLog/clean", controller.CleanCmdbSqlLog)       // 清空SQL操作日志
	// 数据库管理
	router.POST("/cmdb/database", controller.NewCmdbSQLController().CreateDatabase)           // 创建数据库
	router.PUT("/cmdb/database", controller.NewCmdbSQLController().UpdateDatabase)            // 更新数据库
	router.DELETE("/cmdb/database", controller.NewCmdbSQLController().DeleteDatabase)         // 删除数据库
	router.GET("/cmdb/database/info", controller.NewCmdbSQLController().GetDatabase)          // 根据ID获取数据库详情
	router.GET("/cmdb/databaselist", controller.NewCmdbSQLController().ListDatabases)         // 获取数据库列表
	router.GET("/cmdb/database/byname", controller.NewCmdbSQLController().GetDatabasesByName) // 根据名称查询数据库
	router.GET("/cmdb/database/bytype", controller.NewCmdbSQLController().GetDatabasesByType) // 根据类型查询数据库
}
