// utils/system/system.go

package system

import (
	"dodevops-api/api/system/controller"
	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统相关路由
func RegisterSystemRoutes(router *gin.RouterGroup) {
	// 岗位
	router.POST("/post/add", controller.CreateSysPost)
	router.GET("/post/list", controller.GetSysPostList)
	router.GET("/post/info", controller.GetSysPostById)
	router.PUT("/post/update", controller.UpdateSysPost)
	router.DELETE("/post/delete", controller.DeleteSysPostById)
	router.DELETE("/post/batch/delete", controller.BatchDeleteSysPost)
	router.PUT("/post/updateStatus", controller.UpdateSysPostStatus)
	router.GET("/post/vo/list", controller.QuerySysPostVoList)
	//  部门
	router.GET("/dept/list", controller.GetSysDeptList)         // 部门列表
	router.POST("/dept/add", controller.CreateSysDept)          // 添加部门
	router.GET("/dept/info", controller.GetSysDeptById)         // 根据id查询部门
	router.PUT("/dept/update", controller.UpdateSysDept)        // 修改部门
	router.DELETE("/dept/delete", controller.DeleteSysDeptById) // 删除部门
	router.GET("/dept/vo/list", controller.QuerySysDeptVoList)  // 查询部门树形结构
	router.GET("/dept/users", controller.GetDeptUsers)          // 查询部门下的用户
	// 菜单
	router.POST("/menu/add", controller.CreateSysMenu)
	router.GET("/menu/vo/list", controller.QuerySysMenuVoList)
	router.GET("/menu/info", controller.GetSysMenu)
	router.PUT("/menu/update", controller.UpdateSysMenu)
	router.DELETE("/menu/delete", controller.DeleteSysMenu)
	router.GET("/menu/list", controller.GetSysMenuList)
	// 角色
	router.POST("/role/add", controller.CreateSysRole)
	router.GET("/role/info", controller.GetSysRoleById)
	router.PUT("/role/update", controller.UpdateSysRole)
	router.DELETE("/role/delete", controller.DeleteSysRoleById)
	router.PUT("/role/updateStatus", controller.UpdateSysRoleStatus)
	router.GET("/role/list", controller.GetSysRoleList)
	router.GET("/role/vo/list", controller.QuerySysRoleVoList)
	router.GET("/role/vo/idList", controller.QueryRoleMenuIdList)
	router.PUT("/role/assignPermissions", controller.AssignPermissions)
	// 用户
	router.POST("/admin/add", controller.CreateSysAdmin)
	router.GET("/admin/info", controller.GetSysAdminInfo)
	router.PUT("/admin/update", controller.UpdateSysAdmin)
	router.DELETE("/admin/delete", controller.DeleteSysAdminById)
	router.PUT("/admin/updateStatus", controller.UpdateSysAdminStatus)
	router.PUT("/admin/updatePassword", controller.ResetSysAdminPassword)
	router.GET("/admin/list", controller.GetSysAdminList)
	router.POST("/upload", controller.Upload)
	router.PUT("/admin/updatePersonal", controller.UpdatePersonal)
	router.PUT("/admin/updatePersonalPassword", controller.UpdatePersonalPassword)
	// 日志
	router.GET("/sysLoginInfo/list", controller.GetSysLoginInfoList)
	router.DELETE("/sysLoginInfo/batch/delete", controller.BatchDeleteSysLoginInfo)
	router.DELETE("/sysLoginInfo/delete", controller.DeleteSysLoginInfoById)
	router.DELETE("/sysLoginInfo/clean", controller.CleanSysLoginInfo)
	router.GET("/sysOperationLog/list", controller.GetSysOperationLogList)
	router.DELETE("/sysOperationLog/delete", controller.DeleteSysOperationLogById)
	router.DELETE("/sysOperationLog/batch/delete", controller.BatchDeleteSysOperationLog)
	router.DELETE("/sysOperationLog/clean", controller.CleanSysOperationLog)
}
