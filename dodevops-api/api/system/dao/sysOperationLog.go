// 操作日志 数据层
// author xiaoRui

package dao

import (
	"dodevops-api/api/system/model"
	. "dodevops-api/pkg/db"
)

// 新增操作日志
func CreateSysOperationLog(log model.SysOperationLog) {
	Db.Create(&log)
}

// 分页查询操作日志列表
func GetSysOperationLogList(Username, BeginTime, EndTime string, PageSize, PageNum int) (sysOperationLog []model.SysOperationLog, count int64) {
	curDb := Db.Table("sys_operation_log")
	if Username != "" {
		curDb = curDb.Where("username =?", Username)
	}
	if BeginTime != "" && EndTime != "" {
		curDb = curDb.Where("`create_time` BETWEEN ? AND ?", BeginTime, EndTime)
	}
	curDb.Count(&count)
	curDb.Limit(PageSize).Offset((PageNum - 1) * PageSize).Order("create_time desc").Find(&sysOperationLog)
	return sysOperationLog, count
}

// 根据id删除操作日志
func DeleteSysOperationLogById(dto model.SysOperationLogIdDto) {
	Db.Delete(&model.SysOperationLog{}, dto)
}

// 批量删除批量操作日志
func BatchDeleteSysOperationLog(dto model.BatchDeleteSysOperationLogDto) {
	Db.Where("id in (?)", dto.Ids).Delete(&model.SysOperationLog{})
}

// 清空操作日志
func CleanSysOperationLog() {
	Db.Exec("truncate table sys_operation_log")
}
