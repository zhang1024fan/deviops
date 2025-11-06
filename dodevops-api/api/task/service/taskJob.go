package service

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"dodevops-api/api/task/dao"
	"dodevops-api/api/task/model"
	cmdbmodel "dodevops-api/api/cmdb/model"
	"gorm.io/gorm"
	"github.com/robfig/cron/v3"
)

// TaskJobService 任务服务接口
type TaskJobService interface {
	CreateTask(task *model.Task) error
	GetTask(id uint) (*model.Task, error)
	UpdateTask(id uint, task *model.Task) error
	DeleteTask(id uint) error
	ListTasks(page, pageSize int, name string, status int) ([]model.Task, int64, error)
	ListTasksWithDetails(page, pageSize int, name string, status int) ([]model.TaskWithDetails, int64, error)
	TaskNameExists(name string) (bool, error)
	GetTasksByName(name string) ([]model.Task, error)
	GetTasksByType(taskType int) ([]model.Task, error)
	GetTasksByStatus(status int) ([]model.Task, error)
	GetNextExecutionTime(cronExpr string) (time.Time, error)
	GetTaskTemplatesWithStatus(taskId uint) ([]map[string]interface{}, error)
	IncrementExecuteCount(taskId uint) error
	UpdateNextRunTime(taskId uint, nextTime time.Time) error
	UpdateTaskAfterExecution(taskId uint, cronExpr string) error
}

type taskJobServiceImpl struct {
	db     *gorm.DB
	taskDao dao.TaskDao
}

func NewTaskService(db *gorm.DB) TaskJobService {
	return &taskJobServiceImpl{
		db:     db,
		taskDao: dao.NewTaskDao(db),
	}
}

func (s *taskJobServiceImpl) CreateTask(task *model.Task) error {
	// 验证定时任务必须包含CronExpr
	if task.Type == 2 && task.CronExpr == "" {
		return errors.New("定时任务必须包含定时表达式")
	}

	// 创建主任务记录（DAO层会自动创建task_work记录）
	err := s.taskDao.Create(task)
	if err != nil {
		return errors.New("创建任务失败: " + err.Error())
	}

	return nil
}

func (s *taskJobServiceImpl) GetTask(id uint) (*model.Task, error) {
	return s.taskDao.GetById(id)
}

func (s *taskJobServiceImpl) UpdateTask(id uint, task *model.Task) error {
	// 验证定时任务必须包含CronExpr
	if task.Type == 2 && task.CronExpr == "" {
		return errors.New("定时任务必须包含定时表达式")
	}

	existingTask, err := s.taskDao.GetById(id)
	if err != nil {
		return err
	}

	// 更新字段
	existingTask.Name = task.Name
	existingTask.Type = task.Type
	existingTask.Shell = task.Shell
	existingTask.HostIDs = task.HostIDs
	existingTask.CronExpr = task.CronExpr
	existingTask.Status = task.Status
	existingTask.Remark = task.Remark

	return s.taskDao.Update(existingTask)
}

func (s *taskJobServiceImpl) DeleteTask(id uint) error {
	// 先删除关联的子任务
	if err := s.db.Where("task_id = ?", id).Delete(&model.TaskWork{}).Error; err != nil {
		return errors.New("删除子任务失败: " + err.Error())
	}

	// 再删除父任务
	return s.taskDao.Delete(id)
}

func (s *taskJobServiceImpl) ListTasks(page, pageSize int, name string, status int) ([]model.Task, int64, error) {
	offset := (page - 1) * pageSize
	return s.taskDao.List(offset, pageSize, name, status)
}

func (s *taskJobServiceImpl) TaskNameExists(name string) (bool, error) {
	return s.taskDao.ExistsByName(name)
}

func (s *taskJobServiceImpl) GetTasksByName(name string) ([]model.Task, error) {
	return s.taskDao.GetTasksByName(name)
}

func (s *taskJobServiceImpl) GetTasksByType(taskType int) ([]model.Task, error) {
	return s.taskDao.GetTasksByType(taskType)
}

func (s *taskJobServiceImpl) GetTasksByStatus(status int) ([]model.Task, error) {
	return s.taskDao.GetTasksByStatus(status)
}

func (s *taskJobServiceImpl) GetNextExecutionTime(cronExpr string) (time.Time, error) {
	scheduler := cron.New(cron.WithSeconds())
	defer scheduler.Stop()
	
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return time.Time{}, err
	}
	
	return schedule.Next(time.Now()), nil
}

func (s *taskJobServiceImpl) GetTaskTemplatesWithStatus(taskId uint) ([]map[string]interface{}, error) {
	return s.taskDao.GetTaskTemplatesWithStatus(taskId)
}

// IncrementExecuteCount 增加任务执行次数
func (s *taskJobServiceImpl) IncrementExecuteCount(taskId uint) error {
	return s.db.Model(&model.Task{}).Where("id = ?", taskId).
		Update("execute_count", gorm.Expr("execute_count + ?", 1)).Error
}

// UpdateNextRunTime 更新下次执行时间
func (s *taskJobServiceImpl) UpdateNextRunTime(taskId uint, nextTime time.Time) error {
	return s.db.Model(&model.Task{}).Where("id = ?", taskId).
		Update("next_run_time", nextTime).Error
}

// UpdateTaskAfterExecution 任务执行完成后更新execute_count和next_run_time
func (s *taskJobServiceImpl) UpdateTaskAfterExecution(taskId uint, cronExpr string) error {
	// 增加执行次数
	if err := s.IncrementExecuteCount(taskId); err != nil {
		return err
	}

	// 如果是定时任务，更新下次执行时间
	if cronExpr != "" {
		nextTime, err := s.GetNextExecutionTime(cronExpr)
		if err != nil {
			return err
		}
		if err := s.UpdateNextRunTime(taskId, nextTime); err != nil {
			return err
		}
	}

	return nil
}

func (s *taskJobServiceImpl) ListTasksWithDetails(page, pageSize int, name string, status int) ([]model.TaskWithDetails, int64, error) {
	offset := (page - 1) * pageSize

	// 获取基本任务列表
	tasks, total, err := s.taskDao.List(offset, pageSize, name, status)
	if err != nil {
		return nil, 0, err
	}

	// 转换为TaskWithDetails格式
	var tasksWithDetails []model.TaskWithDetails
	for _, task := range tasks {
		taskDetail := model.TaskWithDetails{
			ID:           task.ID,
			Name:         task.Name,
			Type:         task.Type,
			Shell:        task.Shell,
			HostIDs:      task.HostIDs,
			CronExpr:     task.CronExpr,
			Tasklog:      task.Tasklog,
			Status:       task.Status,
			Duration:     task.Duration,
			TaskCount:    task.TaskCount,
			ExecuteCount: task.ExecuteCount,
			NextRunTime:  task.NextRunTime,
			Remark:       task.Remark,
			StartTime:    task.StartTime,
			EndTime:      task.EndTime,
			CreatedAt:    task.CreatedAt,
		}

		// 获取关联的模板信息
		if task.Shell != "" {
			templates, err := s.getTemplatesByIds(task.Shell)
			if err == nil {
				taskDetail.Templates = templates
			}
		}

		// 获取关联的主机信息
		if task.HostIDs != "" {
			hosts, err := s.getHostsByIds(task.HostIDs)
			if err == nil {
				taskDetail.Hosts = hosts
			}
		}

		tasksWithDetails = append(tasksWithDetails, taskDetail)
	}

	return tasksWithDetails, total, nil
}

// getTemplatesByIds 根据模板ID字符串获取模板信息
func (s *taskJobServiceImpl) getTemplatesByIds(shellIds string) ([]model.TaskTemplateInfo, error) {
	if shellIds == "" {
		return []model.TaskTemplateInfo{}, nil
	}

	// 解析逗号分隔的ID字符串
	idStrs := strings.Split(shellIds, ",")
	var ids []uint
	for _, idStr := range idStrs {
		if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) == 0 {
		return []model.TaskTemplateInfo{}, nil
	}

	// 查询模板信息
	var templates []model.TaskTemplate
	if err := s.db.Where("id IN ?", ids).Find(&templates).Error; err != nil {
		return nil, err
	}

	// 转换为TaskTemplateInfo格式
	var templateInfos []model.TaskTemplateInfo
	for _, template := range templates {
		templateInfos = append(templateInfos, model.TaskTemplateInfo{
			ID:      template.ID,
			Name:    template.Name,
			Type:    template.Type,
			Content: template.Content,
		})
	}

	return templateInfos, nil
}

// getHostsByIds 根据主机ID字符串获取主机信息
func (s *taskJobServiceImpl) getHostsByIds(hostIds string) ([]model.HostInfo, error) {
	if hostIds == "" {
		return []model.HostInfo{}, nil
	}

	// 解析逗号分隔的ID字符串
	idStrs := strings.Split(hostIds, ",")
	var ids []uint
	for _, idStr := range idStrs {
		if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}

	if len(ids) == 0 {
		return []model.HostInfo{}, nil
	}

	// 查询主机信息
	var hosts []cmdbmodel.CmdbHost
	if err := s.db.Where("id IN ?", ids).Find(&hosts).Error; err != nil {
		return nil, err
	}

	// 转换为HostInfo格式
	var hostInfos []model.HostInfo
	for _, host := range hosts {
		hostInfos = append(hostInfos, model.HostInfo{
			ID:        host.ID,
			HostName:  host.HostName,
			PrivateIP: host.PrivateIP,
			PublicIP:  host.PublicIP,
			SSHIP:     host.SSHIP,
			Status:    host.Status,
		})
	}

	return hostInfos, nil
}
