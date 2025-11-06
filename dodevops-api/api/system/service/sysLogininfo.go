// 登录日志 服务层
// author xiaoRui

package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"github.com/gin-gonic/gin"
)

type ISysLoginInfoService interface {
	GetSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string, PageSize, pageNum int)
	BatchDeleteSysLoginInfo(c *gin.Context, dto model.DelSysLoginInfoDto)
	DeleteSysLoginInfo(c *gin.Context, dto model.SysLoginInfoIdDto)
	CleanSysLoginInfo(c *gin.Context)
}
type SysLoginInfoServiceImpl struct{}

// 分页获取登录日志列表
func (s SysLoginInfoServiceImpl) GetSysLoginInfoList(c *gin.Context, Username, LoginStatus, BeginTime, EndTime string, PageSize, pageNum int) {
	if PageSize < 1 {
		PageSize = 10
	}
	if pageNum < 1 {
		pageNum = 1
	}
	sysLoginInfo, count := dao.GetSysLoginInfoList(Username, LoginStatus, BeginTime, EndTime, PageSize, pageNum)
	result.Success(c, map[string]interface{}{"total": count, "pageSize": PageSize, "pageNum": pageNum, "list": sysLoginInfo})
}

// 清空登录日志
func (s SysLoginInfoServiceImpl) CleanSysLoginInfo(c *gin.Context) {
	dao.CleanSysLoginInfo()
	result.Success(c, true)
}

// 批量删除登录日志
func (s SysLoginInfoServiceImpl) BatchDeleteSysLoginInfo(c *gin.Context, dto model.DelSysLoginInfoDto) {
	dao.BatchDeleteSysLoginInfo(dto)
	result.Success(c, true)
}

// 根据id删除登录日志
func (s SysLoginInfoServiceImpl) DeleteSysLoginInfo(c *gin.Context, dto model.SysLoginInfoIdDto) {
	dao.DeleteSysLoginInfoById(dto)
	result.Success(c, true)
}

var sysLoginInfoService = SysLoginInfoServiceImpl{}

func SysLoginInfoService() ISysLoginInfoService {
	return &sysLoginInfoService
}
