package dao

import (
	"fmt"
	"dodevops-api/api/task/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

// 缓存项结构
type cacheItem struct {
	data      interface{}
	timestamp time.Time
}

type TaskAnsibleDao struct {
	DB    *gorm.DB
	cache map[string]*cacheItem
	mutex sync.RWMutex
}

func NewTaskAnsibleDao(db *gorm.DB) *TaskAnsibleDao {
	return &TaskAnsibleDao{
		DB:    db,
		cache: make(map[string]*cacheItem),
	}
}

// 缓存相关方法
const cacheTTL = 5 * time.Second // 5秒缓存TTL

func (d *TaskAnsibleDao) getFromCache(key string) (interface{}, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	
	item, exists := d.cache[key]
	if !exists {
		return nil, false
	}
	
	// 检查是否过期
	if time.Since(item.timestamp) > cacheTTL {
		return nil, false
	}
	
	return item.data, true
}

func (d *TaskAnsibleDao) setCache(key string, data interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	d.cache[key] = &cacheItem{
		data:      data,
		timestamp: time.Now(),
	}
}

func (d *TaskAnsibleDao) clearCache(pattern string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	
	// 简单实现：清空所有缓存
	d.cache = make(map[string]*cacheItem)
}

// Create 创建Ansible任务
func (d *TaskAnsibleDao) Create(task *model.TaskAnsible) error {
	return d.DB.Create(task).Error
}

// GetByID 根据ID获取任务
func (d *TaskAnsibleDao) GetByID(id uint) (*model.TaskAnsible, error) {
	var task model.TaskAnsible
	err := d.DB.Where("id = ?", id).First(&task).Error
	return &task, err
}

// Update 更新任务信息
func (d *TaskAnsibleDao) Update(task *model.TaskAnsible) error {
	return d.DB.Save(task).Error
}

// Delete 删除任务（级联删除关联的子任务）
func (d *TaskAnsibleDao) Delete(id uint) error {
	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 首先删除所有关联的子任务
	if err := tx.Where("task_id = ?", id).Delete(&model.TaskAnsibleWork{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 然后删除父任务
	if err := tx.Where("id = ?", id).Delete(&model.TaskAnsible{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// List 获取任务列表
func (d *TaskAnsibleDao) List(page, size int) ([]model.TaskAnsible, int64, error) {
	var tasks []model.TaskAnsible
	var count int64

	err := d.DB.Model(&model.TaskAnsible{}).
		Count(&count).
		Offset((page - 1) * size).
		Limit(size).
		Find(&tasks).Error

	return tasks, count, err
}

// GetByType 根据类型获取任务
func (d *TaskAnsibleDao) GetByType(taskType int) ([]model.TaskAnsible, error) {
	var tasks []model.TaskAnsible
	err := d.DB.Where("type = ?", taskType).Find(&tasks).Error
	return tasks, err
}

// GetByName 根据名称模糊查询任务
func (d *TaskAnsibleDao) GetByName(name string) ([]model.TaskAnsible, error) {
	var tasks []model.TaskAnsible
	err := d.DB.Where("name LIKE ?", "%"+name+"%").Find(&tasks).Error
	return tasks, err
}

// GetWorkByID 根据ID获取子任务 (优化版本：只查询必要字段 + 缓存)
func (d *TaskAnsibleDao) GetWorkByID(taskID, workID uint) (*model.TaskAnsibleWork, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("work_%d_%d", taskID, workID)
	if cached, found := d.getFromCache(cacheKey); found {
		if work, ok := cached.(*model.TaskAnsibleWork); ok {
			return work, nil
		}
	}
	
	var work model.TaskAnsibleWork
	// 只查询必要的字段，减少数据传输
	err := d.DB.Select("id, task_id, entry_file_name, log_path, status, start_time, end_time").
		Where("task_id = ? AND id = ?", taskID, workID).First(&work).Error
	if err != nil {
		return nil, err
	}
	
	// 存入缓存
	d.setCache(cacheKey, &work)
	return &work, nil
}

// GetJobStatus 获取任务状态
func (d *TaskAnsibleDao) GetJobStatus(taskID, workID uint) (*model.TaskAnsibleStatus, error) {
	var work model.TaskAnsibleWork
	err := d.DB.Where("task_id = ? AND id = ?", taskID, workID).First(&work).Error
	if err != nil {
		return nil, err
	}

	status := &model.TaskAnsibleStatus{
		Status: work.Status,
		Log:    "", // 日志从文件读取，不再存储在数据库中
	}

	if work.StartTime != nil {
		status.StartTime = *work.StartTime
	}
	if work.EndTime != nil {
		status.EndTime = *work.EndTime
	}

	return status, nil
}

// GetTaskDetail 获取任务详情 (优化版本：减少预加载数据 + 缓存)
func (d *TaskAnsibleDao) GetTaskDetail(taskID uint) (*model.TaskAnsible, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("task_detail_%d", taskID)
	if cached, found := d.getFromCache(cacheKey); found {
		if task, ok := cached.(*model.TaskAnsible); ok {
			return task, nil
		}
	}
	
	var task model.TaskAnsible
	// 只预加载Works的关键字段，减少数据传输
	err := d.DB.Preload("Works", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, task_id, entry_file_name, status, start_time, end_time, duration")
	}).Where("id = ?", taskID).First(&task).Error
	
	if err == nil {
		// 存入缓存
		d.setCache(cacheKey, &task)
	}
	
	return &task, err
}

// GetWorkStatus 仅获取子任务状态 (轻量级查询 + 缓存)
func (d *TaskAnsibleDao) GetWorkStatus(taskID, workID uint) (int, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("work_status_%d_%d", taskID, workID)
	if cached, found := d.getFromCache(cacheKey); found {
		if status, ok := cached.(int); ok {
			return status, nil
		}
	}
	
	var status int
	err := d.DB.Model(&model.TaskAnsibleWork{}).
		Select("status").
		Where("task_id = ? AND id = ?", taskID, workID).
		Scan(&status).Error
		
	if err == nil {
		// 存入缓存
		d.setCache(cacheKey, status)
	}
	
	return status, err
}

// StartJob 启动任务
func (d *TaskAnsibleDao) StartJob(taskID uint) error {
	err := d.DB.Model(&model.TaskAnsible{}).
		Where("id = ?", taskID).
		Update("status", 2).Error // 2表示运行中
	
	if err == nil {
		// 清空相关缓存
		d.clearCache("")
	}
	
	return err
}

// StopJob 停止任务
func (d *TaskAnsibleDao) StopJob(taskID, workID uint) error {
	return d.DB.Model(&model.TaskAnsibleWork{}).
		Where("task_id = ? AND id = ?", taskID, workID).
		Updates(map[string]interface{}{
			"status":   4, // 4表示已停止
			"end_time": gorm.Expr("NOW()"),
		}).Error
}
