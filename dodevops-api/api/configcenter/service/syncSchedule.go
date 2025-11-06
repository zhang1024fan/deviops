package service

import (
	"encoding/json"
	"fmt"
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common/util"
	"github.com/robfig/cron/v3"
)

type SyncScheduleService struct {
	dao *dao.SyncScheduleDao
}

func NewSyncScheduleService() *SyncScheduleService {
	return &SyncScheduleService{
		dao: dao.NewSyncScheduleDao(),
	}
}

// Create 创建定时同步配置
func (s *SyncScheduleService) Create(syncSchedule *model.SyncSchedule) error {
	// 验证cron表达式
	if err := s.validateCronExpr(syncSchedule.CronExpr); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	// 验证云厂商类型JSON格式
	if err := s.validateKeyTypes(syncSchedule.KeyTypes); err != nil {
		return fmt.Errorf("invalid keyTypes format: %v", err)
	}

	// 计算下次执行时间
	nextTime, err := s.calculateNextRunTime(syncSchedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}
	syncSchedule.NextRunTime = nextTime

	// 设置默认状态
	if syncSchedule.Status == 0 {
		syncSchedule.Status = 1
	}

	return s.dao.Create(syncSchedule)
}

// Update 更新定时同步配置
func (s *SyncScheduleService) Update(syncSchedule *model.SyncSchedule) error {
	// 验证cron表达式
	if err := s.validateCronExpr(syncSchedule.CronExpr); err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	// 验证云厂商类型JSON格式
	if err := s.validateKeyTypes(syncSchedule.KeyTypes); err != nil {
		return fmt.Errorf("invalid keyTypes format: %v", err)
	}

	// 重新计算下次执行时间
	nextTime, err := s.calculateNextRunTime(syncSchedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}
	syncSchedule.NextRunTime = nextTime

	return s.dao.Update(syncSchedule)
}

// Delete 删除定时同步配置
func (s *SyncScheduleService) Delete(id uint) error {
	return s.dao.Delete(id)
}

// GetByID 根据ID查询定时同步配置
func (s *SyncScheduleService) GetByID(id uint) (*model.SyncSchedule, error) {
	return s.dao.GetByID(id)
}

// ListWithPage 分页查询定时同步配置
func (s *SyncScheduleService) ListWithPage(page, pageSize int) ([]model.SyncSchedule, int64, error) {
	return s.dao.ListWithPage(page, pageSize)
}

// GetActiveSchedules 获取所有启用的定时同步配置
func (s *SyncScheduleService) GetActiveSchedules() ([]model.SyncSchedule, error) {
	return s.dao.GetByStatus(1) // 状态为1表示启用
}

// UpdateLastRunTime 更新上次执行时间和下次执行时间
func (s *SyncScheduleService) UpdateLastRunTime(id uint, lastRunTime time.Time) error {
	schedule, err := s.dao.GetByID(id)
	if err != nil {
		return err
	}

	// 计算下次执行时间
	nextTime, err := s.calculateNextRunTime(schedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}

	lastTime := util.HTime{Time: lastRunTime}
	return s.dao.UpdateRunTime(id, &lastTime, nextTime)
}

// UpdateLastRunTimeAndLog 更新上次执行时间、下次执行时间和同步日志
func (s *SyncScheduleService) UpdateLastRunTimeAndLog(id uint, lastRunTime time.Time, syncLog string) error {
	schedule, err := s.dao.GetByID(id)
	if err != nil {
		return err
	}

	// 计算下次执行时间
	nextTime, err := s.calculateNextRunTime(schedule.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}

	lastTime := util.HTime{Time: lastRunTime}
	return s.dao.UpdateRunTimeAndLog(id, &lastTime, nextTime, syncLog)
}

// ToggleStatus 切换配置状态（启用/禁用）
func (s *SyncScheduleService) ToggleStatus(id uint, status int) error {
	return s.dao.UpdateStatus(id, status)
}

// validateCronExpr 验证cron表达式是否有效（标准5位格式：分时日月周）
func (s *SyncScheduleService) validateCronExpr(cronExpr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronExpr)
	return err
}

// validateKeyTypes 验证云厂商类型JSON格式
func (s *SyncScheduleService) validateKeyTypes(keyTypes string) error {
	var types []int
	return json.Unmarshal([]byte(keyTypes), &types)
}

// calculateNextRunTime 计算下次执行时间（标准5位格式）
func (s *SyncScheduleService) calculateNextRunTime(cronExpr string) (*util.HTime, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, err
	}

	nextTime := schedule.Next(time.Now())
	hTime := util.HTime{Time: nextTime}
	return &hTime, nil
}

// ParseKeyTypes 解析云厂商类型JSON字符串
func (s *SyncScheduleService) ParseKeyTypes(keyTypes string) ([]int, error) {
	var types []int
	err := json.Unmarshal([]byte(keyTypes), &types)
	return types, err
}