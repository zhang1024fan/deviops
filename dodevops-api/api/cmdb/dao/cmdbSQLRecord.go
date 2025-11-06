package dao

import (
	"dodevops-api/api/cmdb/model"
	"gorm.io/gorm"
)

type CmdbSQLRecordDao struct {
	db *gorm.DB
}

func NewCmdbSQLRecordDao(db *gorm.DB) *CmdbSQLRecordDao {
	return &CmdbSQLRecordDao{db: db}
}

// CreateRecord 创建SQL执行记录
func (d *CmdbSQLRecordDao) CreateRecord(record *model.CmdbSQLRecord) error {
	return d.db.Create(record).Error
}

// GetRecordsByDatabase 根据数据库名查询记录
func (d *CmdbSQLRecordDao) GetRecordsByDatabase(database string, limit int) ([]model.CmdbSQLRecord, error) {
	var records []model.CmdbSQLRecord
	err := d.db.Where("database = ?", database).
		Order("query_time DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetRecordsByUser 根据执行用户查询记录
func (d *CmdbSQLRecordDao) GetRecordsByUser(username string, limit int) ([]model.CmdbSQLRecord, error) {
	var records []model.CmdbSQLRecord
	err := d.db.Where("exec_user = ?", username).
		Order("query_time DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetRecentRecords 获取最近的执行记录
func (d *CmdbSQLRecordDao) GetRecentRecords(limit int) ([]model.CmdbSQLRecord, error) {
	var records []model.CmdbSQLRecord
	err := d.db.Order("query_time DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetRecordsByPage 分页获取SQL日志
func (d *CmdbSQLRecordDao) GetRecordsByPage(execUser, beginTime, endTime string, pageSize, pageNum int) ([]model.CmdbSQLRecord, int64, error) {
	var records []model.CmdbSQLRecord
	var total int64
	
	query := d.db.Model(&model.CmdbSQLRecord{})
	if execUser != "" {
		query = query.Where("exec_user = ?", execUser)
	}
	if beginTime != "" {
		query = query.Where("query_time >= ?", beginTime)
	}
	if endTime != "" {
		query = query.Where("query_time <= ?", endTime)
	}
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	offset := (pageNum - 1) * pageSize
	err = query.Order("query_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error
		
	return records, total, err
}

// DeleteById 根据ID删除日志
func (d *CmdbSQLRecordDao) DeleteById(id uint) error {
	return d.db.Where("id = ?", id).Delete(&model.CmdbSQLRecord{}).Error
}

// CleanAll 清空所有日志
func (d *CmdbSQLRecordDao) CleanAll() error {
	return d.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.CmdbSQLRecord{}).Error
}
