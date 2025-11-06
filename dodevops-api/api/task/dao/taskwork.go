package dao

import (
	"dodevops-api/api/task/model"
	"dodevops-api/common"
	"time"

	"gorm.io/gorm"
)

type TaskWorkDaoInterface interface {
	Create(job *model.TaskWork) error
	Update(job *model.TaskWork) error
	Delete(id uint) error
	GetByID(id uint) (*model.TaskWork, error)
	GetByTaskID(taskID uint) ([]model.TaskWork, error)
	GetByTaskAndTemplateID(taskID, templateID uint) (*model.TaskWork, error)
	GetParentTask(taskID uint) (*model.Task, error)
	GetTaskJobByID(taskID uint) (*model.Task, error)
	UpdateStatus(id uint, status int) error
	UpdateTiming(id uint, startTime, endTime *time.Time, duration int) error
	UpdateLog(id uint, log string) error
	UpdateLogPath(id uint, logPath string) error
	UpdateScheduledJob(id uint, status int, scheduledTime *time.Time) error
}

type TaskWorkDaoImpl struct {
	db *gorm.DB
}

func (d *TaskWorkDaoImpl) Create(job *model.TaskWork) error {
	return d.db.Create(job).Error
}

func (d *TaskWorkDaoImpl) Update(job *model.TaskWork) error {
	return d.db.Save(job).Error
}

func (d *TaskWorkDaoImpl) Delete(id uint) error {
	return d.db.Delete(&model.TaskWork{}, id).Error
}

func (d *TaskWorkDaoImpl) GetByID(id uint) (*model.TaskWork, error) {
	var job model.TaskWork
	err := d.db.First(&job, id).Error
	return &job, err
}

func (d *TaskWorkDaoImpl) GetByTaskID(taskID uint) ([]model.TaskWork, error) {
	var jobs []model.TaskWork
	err := d.db.Where("task_id = ?", taskID).Find(&jobs).Error
	return jobs, err
}

func (d *TaskWorkDaoImpl) GetByTaskAndTemplateID(taskID, templateID uint) (*model.TaskWork, error) {
	var job model.TaskWork
	err := d.db.Where("task_id = ? AND template_id = ?", taskID, templateID).First(&job).Error
	return &job, err
}

func (d *TaskWorkDaoImpl) UpdateStatus(id uint, status int) error {
	return d.db.Model(&model.TaskWork{}).Where("id = ?", id).Update("status", status).Error
}

func (d *TaskWorkDaoImpl) UpdateLog(id uint, log string) error {
	return d.db.Model(&model.TaskWork{}).Where("id = ?", id).Update("log", log).Error
}

func (d *TaskWorkDaoImpl) UpdateLogPath(id uint, logPath string) error {
	return d.db.Model(&model.TaskWork{}).Where("id = ?", id).Update("log_path", logPath).Error
}

func (d *TaskWorkDaoImpl) UpdateTiming(id uint, startTime, endTime *time.Time, duration int) error {
	updates := map[string]interface{}{
		"start_time": startTime,
		"end_time":   endTime,
		"duration":   duration,
	}
	return d.db.Model(&model.TaskWork{}).Where("id = ?", id).Updates(updates).Error
}

func (d *TaskWorkDaoImpl) GetParentTask(taskID uint) (*model.Task, error) {
	var task model.Task
	err := d.db.Where("id = ?", taskID).First(&task).Error
	return &task, err
}

func (d *TaskWorkDaoImpl) UpdateScheduledJob(id uint, status int, scheduledTime *time.Time) error {
	updates := map[string]interface{}{
		"status":         status,
		"scheduled_time": scheduledTime,
	}
	return d.db.Model(&model.TaskWork{}).Where("id = ?", id).Updates(updates).Error
}

func (d *TaskWorkDaoImpl) GetTaskJobByID(taskID uint) (*model.Task, error) {
	var task model.Task
	err := d.db.Where("id = ?", taskID).First(&task).Error
	return &task, err
}

func TaskWorkDao() TaskWorkDaoInterface {
	return &TaskWorkDaoImpl{
		db: common.GetDB(),
	}
}
