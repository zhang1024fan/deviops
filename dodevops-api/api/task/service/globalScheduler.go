package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"dodevops-api/api/task/model"
	"dodevops-api/common"

	"github.com/robfig/cron/v3"
)

// GlobalScheduler 全局定时调度器
type GlobalScheduler struct {
	cron     *cron.Cron
	entries  map[uint]cron.EntryID // 任务ID -> EntryID映射
	mutex    sync.RWMutex
	running  bool
}

// ScheduledTask 调度任务信息
type ScheduledTask struct {
	TaskID   uint
	CronExpr string
	Handler  func()
}

var (
	globalScheduler *GlobalScheduler
	schedulerOnce   sync.Once
)

// GetGlobalScheduler 获取全局调度器实例（单例模式）
func GetGlobalScheduler() *GlobalScheduler {
	schedulerOnce.Do(func() {
		globalScheduler = &GlobalScheduler{
			cron:    cron.New(), // 使用标准5位格式（分 时 日 月 周）
			entries: make(map[uint]cron.EntryID),
			running: false,
		}
	})
	return globalScheduler
}

// Start 启动全局调度器
func (gs *GlobalScheduler) Start() error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if gs.running {
		return fmt.Errorf("全局调度器已在运行")
	}

	gs.cron.Start()
	gs.running = true
	log.Println("全局调度器启动成功")
	return nil
}

// Stop 停止全局调度器
func (gs *GlobalScheduler) Stop() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if !gs.running {
		return
	}

	gs.cron.Stop()
	gs.running = false
	log.Println("全局调度器已停止")
}

// AddScheduledTask 添加定时任务
func (gs *GlobalScheduler) AddScheduledTask(taskID uint, cronExpr string, handler func()) error {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	log.Printf("添加定时任务请求 - TaskID=%d, CronExpr=%s, 调度器运行状态=%v", taskID, cronExpr, gs.running)

	if !gs.running {
		return fmt.Errorf("全局调度器未运行")
	}

	// 如果任务已存在，先移除
	if entryID, exists := gs.entries[taskID]; exists {
		gs.cron.Remove(entryID)
		log.Printf("移除已存在的定时任务: TaskID=%d, EntryID=%d", taskID, entryID)
	}

	// 添加新的定时任务
	entryID, err := gs.cron.AddFunc(cronExpr, func() {
		log.Printf("执行定时任务: TaskID=%d, CronExpr=%s", taskID, cronExpr)
		handler()
	})

	if err != nil {
		log.Printf("添加定时任务失败: TaskID=%d, Error=%v", taskID, err)
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	// 保存映射
	gs.entries[taskID] = entryID
	log.Printf("定时任务已成功添加: TaskID=%d, CronExpr=%s, EntryID=%d, 当前已注册任务数=%d",
		taskID, cronExpr, entryID, len(gs.entries))
	return nil
}

// RemoveScheduledTask 移除定时任务
func (gs *GlobalScheduler) RemoveScheduledTask(taskID uint) {
	// 第一步：获取 entryID
	gs.mutex.Lock()
	entryID, exists := gs.entries[taskID]
	gs.mutex.Unlock()

	if !exists {
		return
	}

	// 第二步：从 cron 移除（不持有 gs.mutex）
	gs.cron.Remove(entryID)

	// 第三步：从 entries 删除（重新获取锁）
	gs.mutex.Lock()
	delete(gs.entries, taskID)
	gs.mutex.Unlock()

	log.Printf("定时任务已移除: TaskID=%d", taskID)
}

// UpdateScheduledTask 更新定时任务
func (gs *GlobalScheduler) UpdateScheduledTask(taskID uint, cronExpr string, handler func()) error {
	// 移除旧任务
	gs.RemoveScheduledTask(taskID)

	// 添加新任务
	return gs.AddScheduledTask(taskID, cronExpr, handler)
}

// PauseScheduledTask 暂停定时任务
func (gs *GlobalScheduler) PauseScheduledTask(taskID uint) error {
	// 第一步：获取 entryID（持有锁的时间最短）
	gs.mutex.Lock()
	if !gs.running {
		gs.mutex.Unlock()
		return fmt.Errorf("全局调度器未运行")
	}

	log.Printf("暂停任务请求 - TaskID=%d, 当前已注册任务: %v", taskID, gs.entries)

	entryID, exists := gs.entries[taskID]
	if !exists {
		gs.mutex.Unlock()
		log.Printf("任务不在调度器中: TaskID=%d, 已注册的任务IDs: %v", taskID, getMapKeys(gs.entries))
		return fmt.Errorf("任务不存在: TaskID=%d", taskID)
	}
	gs.mutex.Unlock()

	// 第二步：从 cron 中移除任务（不持有 gs.mutex，避免死锁）
	// cron.Remove() 会等待正在执行的任务完成，这可能需要时间
	log.Printf("正在从调度器移除任务: TaskID=%d, EntryID=%d", taskID, entryID)
	gs.cron.Remove(entryID)
	log.Printf("任务已从调度器移除: TaskID=%d", taskID)

	// 第三步：更新数据库状态（不持有 gs.mutex）
	if err := gs.updateTaskStatus(taskID, model.TaskStatusPaused); err != nil {
		log.Printf("更新任务状态失败: TaskID=%d, Error=%v", taskID, err)
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	// 第四步：从 entries 中删除（重新获取锁）
	gs.mutex.Lock()
	delete(gs.entries, taskID)
	gs.mutex.Unlock()

	log.Printf("定时任务已暂停: TaskID=%d", taskID)
	return nil
}

// ResumeScheduledTask 恢复定时任务
func (gs *GlobalScheduler) ResumeScheduledTask(taskID uint) error {
	// 第一步：检查调度器状态和任务是否已注册（持有锁的时间最短）
	gs.mutex.Lock()
	if !gs.running {
		gs.mutex.Unlock()
		return fmt.Errorf("全局调度器未运行")
	}

	if _, exists := gs.entries[taskID]; exists {
		gs.mutex.Unlock()
		return fmt.Errorf("任务已在运行中: TaskID=%d", taskID)
	}
	gs.mutex.Unlock()

	// 第二步：从数据库获取任务信息（不持有 gs.mutex）
	taskService := NewTaskService(common.GetDB())
	task, err := taskService.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("获取任务信息失败: %v", err)
	}

	// 检查任务状态
	if task.Status != model.TaskStatusPaused {
		return fmt.Errorf("任务状态不是暂停状态: TaskID=%d, Status=%d", taskID, task.Status)
	}

	// 检查是否为定时任务
	if task.Type != model.TaskTypeScheduled {
		return fmt.Errorf("只有定时任务可以恢复: TaskID=%d, Type=%d", taskID, task.Type)
	}

	// 计算下次执行时间
	nextTime, err := gs.calculateNextRunTime(task.CronExpr)
	if err != nil {
		return fmt.Errorf("计算下次执行时间失败: %v", err)
	}

	// 第三步：更新数据库状态（不持有 gs.mutex）
	if err := gs.updateTaskStatusAndNextRunTime(taskID, model.TaskStatusRunning, nextTime); err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	// 第四步：重新添加到调度器（使用统一的回调函数）
	err = gs.AddScheduledTask(taskID, task.CronExpr, createScheduledTaskHandler(taskID))

	if err != nil {
		// 回滚状态
		gs.updateTaskStatus(taskID, model.TaskStatusPaused)
		return fmt.Errorf("重新添加定时任务失败: %v", err)
	}

	log.Printf("定时任务已恢复: TaskID=%d, NextRunTime=%v", taskID, nextTime)
	return nil
}

// updateTaskStatus 更新任务状态
func (gs *GlobalScheduler) updateTaskStatus(taskID uint, status int) error {
	db := common.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	return db.Model(&model.Task{}).Where("id = ?", taskID).Update("status", status).Error
}

// updateTaskStatusAndNextRunTime 更新任务状态和下次执行时间
func (gs *GlobalScheduler) updateTaskStatusAndNextRunTime(taskID uint, status int, nextRunTime time.Time) error {
	db := common.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	updates := map[string]interface{}{
		"status":        status,
		"next_run_time": nextRunTime,
	}

	return db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

// calculateNextRunTime 计算下次执行时间
func (gs *GlobalScheduler) calculateNextRunTime(cronExpr string) (time.Time, error) {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return time.Time{}, err
	}

	return schedule.Next(time.Now()), nil
}

// IsRunning 检查调度器是否运行中
func (gs *GlobalScheduler) IsRunning() bool {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()
	return gs.running
}

// GetEntries 获取所有已注册的任务
func (gs *GlobalScheduler) GetEntries() map[uint]cron.EntryID {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	// 复制一份返回
	entries := make(map[uint]cron.EntryID)
	for taskID, entryID := range gs.entries {
		entries[taskID] = entryID
	}
	return entries
}

// GetStats 获取调度器统计信息
func (gs *GlobalScheduler) GetStats() map[string]interface{} {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	stats := map[string]interface{}{
		"running":        gs.running,
		"total_tasks":    len(gs.entries),
		"active_entries": len(gs.cron.Entries()),
	}

	// 如果调度器运行中，获取更详细的信息
	if gs.running {
		entries := gs.cron.Entries()
		var nextRuns []map[string]interface{}

		for _, entry := range entries {
			nextRuns = append(nextRuns, map[string]interface{}{
				"entry_id": entry.ID,
				"next_run": entry.Next.Format("2006-01-02 15:04:05"),
				"prev_run": entry.Prev.Format("2006-01-02 15:04:05"),
			})
		}

		stats["next_runs"] = nextRuns
	}

	return stats
}

// LoadScheduledTasks 从数据库加载所有活跃的定时任务
func (gs *GlobalScheduler) LoadScheduledTasks() error {
	if !gs.IsRunning() {
		return fmt.Errorf("全局调度器未运行")
	}

	// 获取所有定时任务
	taskService := NewTaskService(common.GetDB())
	tasks, err := taskService.GetTasksByType(model.TaskTypeScheduled)
	if err != nil {
		return fmt.Errorf("获取定时任务失败: %v", err)
	}

	log.Printf("开始加载定时任务，共 %d 个", len(tasks))

	// 统计计数器
	var (
		loadedCount   = 0
		skippedCount  = 0
		resumedCount  = 0
	)

	// 逐个处理定时任务
	for _, task := range tasks {
		taskID := task.ID

		log.Printf("处理定时任务: TaskID=%d, Type=%d, Status=%d, CronExpr=%s",
			taskID, task.Type, task.Status, task.CronExpr)

		// 跳过非活跃状态的任务
		if task.Status == model.TaskStatusFailed {
			log.Printf("跳过失败任务: TaskID=%d, Status=%d", taskID, task.Status)
			skippedCount++
			continue
		}

		// 跳过没有cron表达式的任务
		if task.CronExpr == "" {
			log.Printf("任务 %d 缺少cron表达式，跳过", taskID)
			skippedCount++
			continue
		}

		// 检查任务状态并处理
		switch task.Status {
		case model.TaskStatusRunning:
			// 服务重启后，恢复运行中的定时任务
			log.Printf("恢复运行中的定时任务: TaskID=%d", taskID)

			// 重新计算下次执行时间
			nextTime, err := gs.calculateNextRunTime(task.CronExpr)
			if err != nil {
				log.Printf("计算下次执行时间失败: TaskID=%d, Error=%v", taskID, err)
				skippedCount++
				continue
			}

			// 更新下次执行时间
			if err := gs.updateTaskStatusAndNextRunTime(taskID, model.TaskStatusRunning, nextTime); err != nil {
				log.Printf("更新任务下次执行时间失败: TaskID=%d, Error=%v", taskID, err)
			}

			// 添加到调度器
			if err := gs.addTaskToScheduler(taskID, task.CronExpr); err != nil {
				log.Printf("恢复定时任务失败: TaskID=%d, Error=%v", taskID, err)
				skippedCount++
				continue
			}

			resumedCount++
			loadedCount++

		case model.TaskStatusPending:
			// 等待中的定时任务，启动并添加到调度器
			log.Printf("启动等待中的定时任务: TaskID=%d", taskID)

			// 计算下次执行时间
			nextTime, err := gs.calculateNextRunTime(task.CronExpr)
			if err != nil {
				log.Printf("计算下次执行时间失败: TaskID=%d, Error=%v", taskID, err)
				skippedCount++
				continue
			}

			// 更新状态为运行中并设置下次执行时间
			if err := gs.updateTaskStatusAndNextRunTime(taskID, model.TaskStatusRunning, nextTime); err != nil {
				log.Printf("更新任务状态失败: TaskID=%d, Error=%v", taskID, err)
				skippedCount++
				continue
			}

			// 添加到调度器
			if err := gs.addTaskToScheduler(taskID, task.CronExpr); err != nil {
				log.Printf("启动定时任务失败: TaskID=%d, Error=%v", taskID, err)
				// 回滚状态
				gs.updateTaskStatus(taskID, model.TaskStatusPending)
				skippedCount++
				continue
			}

			loadedCount++

		case model.TaskStatusPaused:
			// 暂停状态的任务不需要添加到调度器
			log.Printf("跳过暂停状态任务: TaskID=%d", taskID)
			skippedCount++

		default:
			log.Printf("跳过其他状态任务: TaskID=%d, Status=%d", taskID, task.Status)
			skippedCount++
		}
	}

	log.Printf("定时任务加载完成 - 总数: %d, 已加载: %d, 已恢复: %d, 已跳过: %d",
		len(tasks), loadedCount, resumedCount, skippedCount)

	return nil
}

// createScheduledTaskHandler 创建统一的定时任务回调函数
func createScheduledTaskHandler(taskID uint) func() {
	return func() {
		log.Printf("\n=== 定时任务触发 (LoadScheduledTasks) ===")
		log.Printf("触发时间: %v", time.Now().Format("2006-01-02 15:04:05"))
		log.Printf("任务ID: %d", taskID)
		log.Println("----------------------------")

		// 获取任务服务
		workService := NewTaskWorkService()
		workDao := workService.(*TaskWorkServiceImpl).dao

		// 执行前检查任务状态
		taskService := NewTaskService(common.GetDB())
		currentTask, err := taskService.GetTask(taskID)
		if err != nil {
			log.Printf("获取任务失败: TaskID=%d, Error=%v", taskID, err)
			return
		}
		log.Printf("父任务状态: %d", currentTask.Status)

		// 检查任务是否被暂停或停止
		if currentTask.Status == model.TaskStatusPaused {
			log.Printf("定时任务已被暂停，跳过执行: TaskID=%d", taskID)
			return
		}

		if currentTask.Status == model.TaskStatusFailed {
			log.Printf("定时任务状态异常，跳过执行: TaskID=%d", taskID)
			return
		}

		// 获取任务队列服务
		taskQueue := GetTaskQueue()
		if taskQueue == nil {
			log.Printf("任务队列服务未初始化，跳过执行: TaskID=%d", taskID)
			return
		}
		log.Printf("任务队列服务已就绪")

		// 获取所有子任务
		currentJobs, err := workDao.GetByTaskID(taskID)
		if err != nil {
			log.Printf("获取子任务失败: TaskID=%d, Error=%v", taskID, err)
			return
		}
		log.Printf("获取到 %d 个子任务", len(currentJobs))

		// 提交所有子任务到队列
		submittedCount := 0
		for _, currentJob := range currentJobs {
			// 跳过已停止或已完成的任务
			if currentJob.Status == 4 { // 异常
				log.Printf("跳过异常状态的子任务: ID=%d, Status=%d", currentJob.ID, currentJob.Status)
				continue
			}
			if currentJob.Status == 3 { // 已完成
				log.Printf("跳过已完成的子任务: ID=%d, Status=%d", currentJob.ID, currentJob.Status)
				continue
			}
			if currentJob.Status == 2 { // 运行中
				log.Printf("子任务正在运行中，跳过: ID=%d", currentJob.ID)
				continue
			}

			log.Printf("准备执行子任务: ID=%d, TemplateID=%d, 当前状态=%d",
				currentJob.ID, currentJob.TemplateID, currentJob.Status)

			// 更新状态为运行中
			if err := workDao.UpdateStatus(currentJob.ID, 2); err != nil {
				log.Printf("更新子任务状态失败: %v", err)
				continue
			}

			// 提交到队列
			priority := workService.(*TaskWorkServiceImpl).getTaskPriority(&currentJob)
			if err := taskQueue.Enqueue(&currentJob, priority); err != nil {
				log.Printf("提交定时任务到队列失败: %v", err)
				// 回滚状态
				workDao.UpdateStatus(currentJob.ID, 1)
				continue
			}

			submittedCount++
			log.Printf("定时任务已提交到队列: TaskID=%d, JobID=%d, Priority=%s",
				currentJob.TaskID, currentJob.ID, priority)
		}

		log.Printf("共提交 %d 个子任务到队列", submittedCount)
		log.Println("===================")
	}
}

// addTaskToScheduler 添加任务到调度器的通用方法
func (gs *GlobalScheduler) addTaskToScheduler(taskID uint, cronExpr string) error {
	return gs.AddScheduledTask(taskID, cronExpr, createScheduledTaskHandler(taskID))
}

// 全局函数用于初始化和管理
func InitGlobalScheduler() error {
	scheduler := GetGlobalScheduler()
	if err := scheduler.Start(); err != nil {
		return err
	}

	// 加载已有的定时任务
	if err := scheduler.LoadScheduledTasks(); err != nil {
		log.Printf("加载定时任务失败: %v", err)
		// 不返回错误，允许系统继续运行
	}

	return nil
}

func StopGlobalScheduler() {
	if globalScheduler != nil {
		globalScheduler.Stop()
	}
}

// getMapKeys 获取map的所有key
func getMapKeys(m map[uint]cron.EntryID) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}