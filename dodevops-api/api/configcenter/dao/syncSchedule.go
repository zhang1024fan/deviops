package dao

import (
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common"
	"dodevops-api/common/util"
	"gorm.io/gorm"
)

type SyncScheduleDao struct {
	db *gorm.DB
}

func NewSyncScheduleDao() *SyncScheduleDao {
	return &SyncScheduleDao{
		db: common.GetDB(),
	}
}

// Create 创建定时同步配置
func (d *SyncScheduleDao) Create(syncSchedule *model.SyncSchedule) error {
	return d.db.Create(syncSchedule).Error
}

// Update 更新定时同步配置
func (d *SyncScheduleDao) Update(syncSchedule *model.SyncSchedule) error {
	return d.db.Save(syncSchedule).Error
}

// Delete 删除定时同步配置
func (d *SyncScheduleDao) Delete(id uint) error {
	return d.db.Delete(&model.SyncSchedule{}, id).Error
}

// GetByID 根据ID查询定时同步配置
func (d *SyncScheduleDao) GetByID(id uint) (*model.SyncSchedule, error) {
	var syncSchedule model.SyncSchedule
	err := d.db.First(&syncSchedule, id).Error
	if err != nil {
		return nil, err
	}
	return &syncSchedule, nil
}

// ListWithPage 分页查询定时同步配置
func (d *SyncScheduleDao) ListWithPage(page, pageSize int) ([]model.SyncSchedule, int64, error) {
	var syncSchedules []model.SyncSchedule
	var total int64

	// 计算分页
	offset := (page - 1) * pageSize

	// 查询总数
	if err := d.db.Model(&model.SyncSchedule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	err := d.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&syncSchedules).Error
	if err != nil {
		return nil, 0, err
	}

	return syncSchedules, total, nil
}

// GetByStatus 根据状态查询定时同步配置
func (d *SyncScheduleDao) GetByStatus(status int) ([]model.SyncSchedule, error) {
	var syncSchedules []model.SyncSchedule
	err := d.db.Where("status = ?", status).Find(&syncSchedules).Error
	return syncSchedules, err
}

// UpdateStatus 更新配置状态
func (d *SyncScheduleDao) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRunTime 更新执行时间
func (d *SyncScheduleDao) UpdateRunTime(id uint, lastRunTime, nextRunTime *util.HTime) error {
	updates := map[string]interface{}{
		"last_run_time": lastRunTime,
		"next_run_time": nextRunTime,
	}
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateRunTimeAndLog 更新执行时间和同步日志
func (d *SyncScheduleDao) UpdateRunTimeAndLog(id uint, lastRunTime, nextRunTime *util.HTime, syncLog string) error {
	updates := map[string]interface{}{
		"last_run_time": lastRunTime,
		"next_run_time": nextRunTime,
		"sync_log":      syncLog,
	}
	return d.db.Model(&model.SyncSchedule{}).Where("id = ?", id).Updates(updates).Error
}

// GetAll 获取所有定时同步配置
func (d *SyncScheduleDao) GetAll() ([]model.SyncSchedule, error) {
	var syncSchedules []model.SyncSchedule
	err := d.db.Find(&syncSchedules).Error
	return syncSchedules, err
}