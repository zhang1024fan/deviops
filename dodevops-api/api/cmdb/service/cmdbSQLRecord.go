package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"time"
	"github.com/gin-gonic/gin"
)

type CmdbSQLRecordService struct {
	recordDao *dao.CmdbSQLRecordDao
}

func NewCmdbSQLRecordService(recordDao *dao.CmdbSQLRecordDao) *CmdbSQLRecordService {
	return &CmdbSQLRecordService{recordDao: recordDao}
}

// RecordSQLExecution 记录SQL执行
func (s *CmdbSQLRecordService) RecordSQLExecution(
	instanceID string,
	database string,
	operationType string,
	sqlContent string,
	execUser string,
	ip string,
	scannedRows int64,
	affectedRows int64,
	executionTime int64,
	returnedRows int64,
	result string,
) error {
	record := &model.CmdbSQLRecord{
		InstanceID:    instanceID,
		Database:      database,
		OperationType: operationType,
		SQLContent:    sqlContent,
		ExecUser:      execUser,
		IP:            ip,
		ScannedRows:   scannedRows,
		AffectedRows:  affectedRows,
		ExecutionTime: executionTime,
		ReturnedRows:  returnedRows,
		Result:        result,
		QueryTime:     util.HTime{Time: time.Now()},
	}

	return s.recordDao.CreateRecord(record)
}

// GetRecentSQLRecords 获取最近的SQL执行记录
func (s *CmdbSQLRecordService) GetRecentSQLRecords(limit int) ([]model.CmdbSQLRecord, error) {
	return s.recordDao.GetRecentRecords(limit)
}

// GetSQLRecordsByDatabase 根据数据库获取SQL执行记录
func (s *CmdbSQLRecordService) GetSQLRecordsByDatabase(database string, limit int) ([]model.CmdbSQLRecord, error) {
	return s.recordDao.GetRecordsByDatabase(database, limit)
}

// GetSQLRecordsByUser 根据用户获取SQL执行记录
func (s *CmdbSQLRecordService) GetSQLRecordsByUser(username string, limit int) ([]model.CmdbSQLRecord, error) {
	return s.recordDao.GetRecordsByUser(username, limit)
}

// GetCmdbSqlLogList 分页获取SQL日志
func (s *CmdbSQLRecordService) GetCmdbSqlLogList(c *gin.Context, execUser, beginTime, endTime string, pageSize, pageNum int) {
	records, total, err := s.recordDao.GetRecordsByPage(execUser, beginTime, endTime, pageSize, pageNum)
	if err != nil {
		result.Failed(c, 500, err.Error())
		return
	}
	result.Success(c, gin.H{
		"records": records,
		"total":   total,
	})
}

// DeleteCmdbSqlLogById 根据ID删除日志
func (s *CmdbSQLRecordService) DeleteCmdbSqlLogById(c *gin.Context, dto model.CmdbSqlLogIdDto) {
	if err := s.recordDao.DeleteById(dto.Id); err != nil {
		result.Failed(c, 500, err.Error())
		return
	}
	result.Success(c, nil)
}

// CleanCmdbSqlLog 清空SQL日志
func (s *CmdbSQLRecordService) CleanCmdbSqlLog(c *gin.Context) {
	if err := s.recordDao.CleanAll(); err != nil {
		result.Failed(c, 500, err.Error())
		return
	}
	result.Success(c, nil)
}
